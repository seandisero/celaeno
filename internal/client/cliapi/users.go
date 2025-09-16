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

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode == 409 {
		return shared.User{}, fmt.Errorf("user %s already exists", name)
	}
	if resp.StatusCode > 299 {
		return shared.User{}, fmt.Errorf("error while registering user %v\n status code: %d", name, resp.StatusCode)
	}

	var user shared.User
	err = decoder.Decode(&user)
	if err != nil {
		return shared.User{}, fmt.Errorf("error while decoding new user data: %w", err)
	}

	return user, nil
}

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

	err = auth.SaveTokenToFile(loginData.Token)
	if err != nil {
		return shared.User{}, err
	}

	user := shared.User{
		ID:        loginData.User.ID,
		Username:  loginData.User.Username,
		CreatedAt: loginData.User.CreatedAt,
		UpdatedAt: loginData.User.UpdatedAt,
	}

	return user, nil

}

func (cli *CelaenoClient) PostMessage(message string) error {
	reqBody := shared.Message{
		Message: message,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	reader := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", "http://localhost:8080/app", reader)
	if err != nil {
		return fmt.Errorf("failed to create new request: %v", err)
	}

	authToken, err := auth.AuthToken()
	if err != nil {
		return fmt.Errorf("could not get auth token: %w", err)
	}

	bearerToken := "Bearer " + authToken
	req.Header.Set("Authorization", bearerToken)

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("something went wrong do-ing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var respError shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respError)
		if err != nil {
			return fmt.Errorf("error decoding error responce %w", err)
		}
		return fmt.Errorf("error returned from server: %s\nwith status code: %d", respError.Error, resp.StatusCode)
	}

	var respBody shared.Message
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		slog.Error("this is where my error is")
		slog.Info(respBody.Message)
		return fmt.Errorf("error unmarshaling json body")
	}
	fmt.Printf("\t\t\t%s < \n", respBody.Message)
	return nil

}
