package cliapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/client/screen"
	"github.com/seandisero/celaeno/internal/shared"
)

const (
	UserAgent = "celaeno-client"
)

type CelaenoClient struct {
	HttpClient *http.Client
	Connection *websocket.Conn
	Cancel     context.CancelFunc
	Ctx        context.Context
	WS_URL     string
	URL        string
	LocalUser  *shared.User
	ChatRoom   string
	Screen     *screen.Screen
}

func NewClient(timeout time.Duration) *CelaenoClient {
	client := CelaenoClient{
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return &client
}

func (cli *CelaenoClient) readConnection() (shared.Message, error) {
	var message shared.Message
	_, msg, err := cli.Connection.Reader(context.Background())
	if err != nil {
		if err == io.EOF {
			return message, err
		} else {
			return message, nil
		}
	}

	err = json.NewDecoder(msg).Decode(&message)
	if err != nil {
		return message, fmt.Errorf("Could not decode message %s", err.Error())
	}

	message.Incoming = true
	if message.Username == cli.LocalUser.Username {
		message.Incoming = false
	}

	return message, nil
}

func (cli *CelaenoClient) Listen() {
	for {
		if <-cli.Ctx.Done() {
			return
		}
		_, msg, err := cli.Connection.Reader(context.Background())
		if err != nil {
			if err == io.EOF {
				break
			} else {
				continue
			}
		}

		var message shared.Message
		err = json.NewDecoder(msg).Decode(&message)
		if err != nil {
			fmt.Printf("Could not decode message %s", err.Error())
			continue
		}
		if message.Message == "/exit" {
			cli.Cancel()
			break
		}

		message.Incoming = true
		if message.Username == cli.LocalUser.Username {
			message.Incoming = false
		}
		cli.Screen.ClearMessageBox()
		cli.Screen.HandleMessage(message)
	}
	slog.Warn("client is no longer listening")
}
