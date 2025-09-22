package cliapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

func (cli *CelaenoClient) Connect(name string) error {
	if cli.Connection != nil {
		err := cli.Connection.CloseNow()
		if err != nil {
			slog.Error("error clearing input", "error", err)
			return err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())

	if cli.ConnCtx != nil {
		cli.Cancel()
	}
	cli.Cancel = cancel
	cli.ConnCtx = ctx

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

	connectionURL := cli.URL + "/api/chat/connect/" + name
	conn, resp, err := websocket.Dial(ctx, connectionURL, &options)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		if err != nil {
			slog.Error("json could not decode responce error", "error", err)
		}
		return fmt.Errorf("%s", respErr.Error)
	}

	cli.Connection = conn
	go cli.Listen()

	return nil
}
