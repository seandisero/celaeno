package cliapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

func (cli *CelaenoClient) Login(name, password string) (shared.User, *websocket.Conn, error) {
	loginRequest := shared.LoginRequest{
		Name:     name,
		Password: password,
	}

	jsonBody, err := json.Marshal(loginRequest)
	if err != nil {
		return shared.User{}, nil, fmt.Errorf("could not marshal request")
	}

	reqBody := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", cli.URL+"/api/login", reqBody)
	if err != nil {
		return shared.User{}, nil, fmt.Errorf("error creating new request: %w", err)
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return shared.User{}, nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		return shared.User{}, nil, fmt.Errorf("%s", respErr.Error)
	}

	decoder := json.NewDecoder(resp.Body)
	type loginResponce struct {
		User  shared.User
		Token string `json:"token"`
	}

	var loginData loginResponce
	err = decoder.Decode(&loginData)
	if err != nil {
		return shared.User{}, nil, fmt.Errorf("could not decode responce data: %w", err)
	}

	err = auth.SetAuthToken(loginData.Token)
	if err != nil {
		return shared.User{}, nil, err
	}

	user := shared.User{
		ID:          loginData.User.ID,
		Username:    loginData.User.Username,
		Displayname: loginData.User.Displayname,
		CreatedAt:   loginData.User.CreatedAt,
		UpdatedAt:   loginData.User.UpdatedAt,
	}

	cli.LocalUser = &user

	conn, err := cli.establishConnection()
	if err != nil {
		return user, nil, fmt.Errorf("Could not establish websocket connection: %w", err)
	}

	if cli.Connection == nil {
		slog.Info("", "connection", "nil")
	}

	return user, conn, nil
}

func (cli *CelaenoClient) establishConnection() (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	if cli.Cancel != nil {
		cli.Cancel()
	}
	cli.Cancel = cancel

	header := http.Header{}
	token, err := auth.AuthToken()
	if err != nil {
		return nil, fmt.Errorf("error retireving auth token: %w", err)
	}
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", UserAgent)

	options := websocket.DialOptions{
		HTTPClient: cli.HttpClient,
		HTTPHeader: header,
	}

	conn, resp, err := websocket.Dial(ctx, cli.WS_URL+"/api/chat/ws", &options)
	if err != nil {
		return nil, err
	}
	slog.Info("no error from dial")
	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		return nil, fmt.Errorf("%s", respErr.Error)
	}
	slog.Info("no responce error")

	cli.Connection = conn
	go cli.Listen()

	return conn, nil
}

func (cli *CelaenoClient) Logout() error {

	err := auth.SetAuthToken("")
	if err != nil {
		return fmt.Errorf("could not remove token on logout: %w", err)
	}

	return nil
}
