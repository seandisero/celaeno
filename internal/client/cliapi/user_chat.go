package cliapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

func (cli *CelaenoClient) PostMessage(message *shared.Message) error {
	if cli.Connection == nil {
		return fmt.Errorf("you must make a connection or start a chat to post a message")
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal data")
	}

	data := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", cli.URL+"/api/chat/publish/"+cli.ChatRoom, data)
	if err != nil {
		return fmt.Errorf("unsucessful creation of request: %w", err)
	}

	err = auth.ApplyBearerToken(req, cli.LocalUser.Username)
	if err != nil {
		return err
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		slog.Error("unsucessful request", "error", err)
		return fmt.Errorf("unsucessful request: %w", err)
	}

	if resp.StatusCode > 299 {
		cli.Connection.Close(websocket.StatusInternalError, "host has left the chat")
		return fmt.Errorf("chat no longer exists")
	}

	return nil
}
