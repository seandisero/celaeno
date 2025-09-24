package cliapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/seandisero/celaeno/internal/client/auth"
)

func (cli CelaenoClient) ListChatRooms() ([]string, error) {
	type responceBody struct {
		ChatRooms []string `json:"chat_rooms"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", cli.URL+"/api/chat/rooms", http.NoBody)
	if err != nil {
		return []string{}, err
	}

	err = auth.ApplyBearerToken(req, cli.LocalUser.Username)
	if err != nil {
		return []string{}, err
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	var respBody responceBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return []string{}, err
	}
	return respBody.ChatRooms, nil
}
