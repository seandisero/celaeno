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

func (cli *CelaenoClient) RegisterUser(name, password string) (shared.User, error) {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	requestBody := request{
		Name:     name,
		Password: password,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return shared.User{}, fmt.Errorf("error marshaling json data")
	}

	reader := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", cli.URL+"/api/users", reader)
	if err != nil {
		return shared.User{}, fmt.Errorf("error creating new request")
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return shared.User{}, fmt.Errorf("error while doing request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 409 {
		return shared.User{}, fmt.Errorf("user %s already exists", name)
	}
	if resp.StatusCode > 299 {
		return shared.User{}, fmt.Errorf("error while registering user %v\n status code: %d", name, resp.StatusCode)
	}

	var user shared.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return shared.User{}, fmt.Errorf("error while decoding new user data: %w", err)
	}

	return user, nil
}

func (cli *CelaenoClient) DeleteUser(password string) error {
	type deleteRequest struct {
		Password string `json:"password"`
	}
	reqBody := deleteRequest{
		Password: password,
	}

	user, err := cli.GetUser()
	if err != nil {
		return fmt.Errorf("user not logged in: %w", err)
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshaling request body")
	}

	requestBuffer := bytes.NewBuffer(jsonReq)
	req, err := http.NewRequest("DELETE", cli.URL+fmt.Sprintf("/api/users/{%s}", user.ID), requestBuffer)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	err = auth.ApplyBearerToken(req, user.Username)
	if err != nil {
		return err
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request")
	}

	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		if err != nil {
			slog.Error("json could not decode responce error", "error", err)
		}
		return fmt.Errorf("%s", respErr.Error)
	}

	err = auth.SetAuthToken("", user.Username)
	if err != nil {
		return err
	}

	return nil
}
