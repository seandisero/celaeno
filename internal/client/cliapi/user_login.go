package cliapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

func (cli *CelaenoClient) Login(name, password string) (shared.User, error) {
	loginRequest := shared.LoginRequest{
		Name:     name,
		Password: password,
	}

	jsonBody, err := json.Marshal(loginRequest)
	if err != nil {
		return shared.User{}, fmt.Errorf("could not marshal request")
	}

	reqBody := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", cli.URL+"/api/login", reqBody)
	if err != nil {
		return shared.User{}, fmt.Errorf("error creating new request: %w", err)
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return shared.User{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		if err != nil {
			slog.Error("json could not decode responce error", "error", err)
		}
		return shared.User{}, fmt.Errorf("%s", respErr.Error)
	}

	decoder := json.NewDecoder(resp.Body)
	type loginResponce struct {
		User  shared.User
		Token string `json:"token"`
	}

	var loginData loginResponce
	err = decoder.Decode(&loginData)
	if err != nil {
		return shared.User{}, fmt.Errorf("could not decode responce data: %w", err)
	}

	err = auth.SetAuthToken(loginData.Token, loginData.User.Username)
	if err != nil {
		return shared.User{}, err
	}

	fmt.Printf("displayname: %s", loginData.User.Displayname.String)

	user := shared.User{
		ID:          loginData.User.ID,
		Username:    loginData.User.Username,
		Displayname: loginData.User.Displayname,
		CreatedAt:   loginData.User.CreatedAt,
		UpdatedAt:   loginData.User.UpdatedAt,
	}

	cli.LocalUser = &user

	return user, nil
}

func (cli *CelaenoClient) Logout() error {

	err := auth.SetAuthToken("", cli.LocalUser.Username)
	if err != nil {
		return fmt.Errorf("could not remove token on logout: %w", err)
	}

	return nil
}
