package cliapi

import (
	"context"
	"crypto/cipher"
	"encoding/base64"
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
	ConnCtx    context.Context
	URL        string
	LocalUser  *shared.User
	ChatRoom   string
	Screen     *screen.Screen
	Block      *cipher.Block
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
			return message, err
		}
	}

	err = json.NewDecoder(msg).Decode(&message)
	if err != nil {
		return message, fmt.Errorf("could not decode message %s", err.Error())
	}

	message.Incoming = true
	if message.Username == cli.LocalUser.Username || message.Username == cli.LocalUser.Displayname.String {
		message.Incoming = false
	}

	return message, nil
}

func (cli *CelaenoClient) Listen() {
outer:
	for {
		select {
		case <-cli.ConnCtx.Done():
			break outer
		default:
			message, err := cli.readConnection()
			if err != nil {
				if err == io.EOF {
					break outer
				} else {
					continue
				}
			}

			messageBytes, err := base64.StdEncoding.DecodeString(message.Message)
			if err != nil {
				fmt.Println("could not decode string: %w", err)
			}
			decryption, err := cli.Decrypt(messageBytes)
			if err != nil {
				fmt.Printf("could not decrypt message: %v", err)
				continue
			}
			message.Message = string(decryption)
			cli.Screen.HandleMessage(message)
		}
	}
	err := cli.Connection.CloseNow()
	if err != nil {
		slog.Error("Failed to close connection", "error", err)
		return
	}
	slog.Warn("client is no longer listening")
}
