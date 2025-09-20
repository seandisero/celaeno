package cliapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

func (cli *CelaenoClient) CreateChat() error {
	if cli.Connection != nil {
		cli.Connection.CloseNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	if cli.Cancel != nil {
		cli.Cancel()
	}
	cli.Cancel = cancel

	header := http.Header{}
	token, err := auth.AuthToken(cli.LocalUser.Username)
	if err != nil {
		return fmt.Errorf("error retireving auth token: %w", err)
	}
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", UserAgent)

	options := websocket.DialOptions{
		HTTPClient: cli.HttpClient,
		HTTPHeader: header,
	}

	conn, resp, err := websocket.Dial(ctx, cli.URL+"/api/chat/ws", &options)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		return fmt.Errorf("%s", respErr.Error)
	}

	cli.Connection = conn
	go cli.Listen()

	return nil
}
