package cliapi

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	auth.ApplyBearerToken(req)

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request")
	}

	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
		return fmt.Errorf("%s", respErr.Error)
	}

	auth.SetAuthToken("")

	return nil
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

	if resp.StatusCode > 299 {
		var respErr shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&respErr)
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

	err = auth.SetAuthToken(loginData.Token)
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

func (cli *CelaenoClient) Logout() error {

	err := auth.SetAuthToken("")
	if err != nil {
		return fmt.Errorf("could not remove token on logout: %w", err)
	}

	return nil
}

func (cli *CelaenoClient) GetUser() (shared.User, error) {
	req, err := http.NewRequest("GET", cli.URL+"/api/login", http.NoBody)
	if err != nil {
		return shared.User{}, fmt.Errorf("error creating new request: %w", err)
	}

	auth.ApplyBearerToken(req)

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return shared.User{}, fmt.Errorf("error requesting user data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var responceError shared.ResponceError
		err = json.NewDecoder(resp.Body).Decode(&responceError)
		if err != nil {
			return shared.User{}, fmt.Errorf("error performing request and decoding message")
		}
		return shared.User{}, fmt.Errorf("error performing request: %s", responceError.Error)
	}

	var user shared.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return shared.User{}, fmt.Errorf("error decoding user data: %w", err)
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

	err = auth.ApplyBearerToken(req)
	if err != nil {
		return fmt.Errorf("could not write auth token to header: %w", err)
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("something went wrong do-ing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var ee shared.ResponceError
		if err = json.NewDecoder(resp.Body).Decode(&ee); err != nil {
			return fmt.Errorf("error decoding error responce %w", err)
		}
		return fmt.Errorf("error returned from server: %s\n", ee.Error)
	}

	var respBody shared.Message
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return fmt.Errorf("error unmarshaling json body")
	}
	for i := 0; i < 4; i++ {
		fmt.Printf("\t")
	}
	fmt.Printf("%s < \n", respBody.Message)
	return nil

}

func (cli *CelaenoClient) SetDisplayName(displayname string) (shared.User, error) {
	type request struct {
		DisplayName string `json:"displayname"`
	}

	requestBody := request{
		DisplayName: displayname,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return shared.User{}, fmt.Errorf("could not marshal request: %w", err)
	}

	requestBuffer := bytes.NewBuffer(jsonBody)

	user, err := cli.GetUser()
	if err != nil {
		return shared.User{}, fmt.Errorf("could not get user")
	}

	requestURL := cli.URL + "/api/users/" + string(user.ID)
	req, err := http.NewRequest("PUT", requestURL, requestBuffer)
	if err != nil {
		return shared.User{}, fmt.Errorf("error creating new request: %w", err)
	}

	auth.ApplyBearerToken(req)

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return shared.User{}, fmt.Errorf("error doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var ee shared.ResponceError
		if err = json.NewDecoder(resp.Body).Decode(&ee); err != nil {
			return shared.User{}, fmt.Errorf("error decoding error responce %w", err)
		}
		return shared.User{}, fmt.Errorf("error returned from server: %s\n", ee.Error)
	}

	var data shared.User
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return shared.User{}, fmt.Errorf("error decoding request body: %w", err)
	}

	return data, nil
}
