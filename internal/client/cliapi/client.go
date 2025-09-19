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
	"github.com/seandisero/celaeno/internal/shared"
)

const (
	UserAgent = "celaeno-client"
)

type CelaenoClient struct {
	HttpClient *http.Client
	Connection *websocket.Conn
	Cancel     context.CancelFunc
	WS_URL     string
	URL        string
	LocalUser  *shared.User
	ChatRoom   string
}

func NewClient(timeout time.Duration) *CelaenoClient {
	client := CelaenoClient{
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return &client
}

func (cli *CelaenoClient) Listen() {
	slog.Info("client is listening for messages")
	defer cli.Cancel()
	for {
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

		fmt.Printf("\t\t\t\t< %s\n", message.Username)
		fmt.Printf("\t\t\t\t%s  \n", message.Message)
	}
	slog.Info("client is no longer listening for messages")
}
