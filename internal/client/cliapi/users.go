package cliapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/shared"
)

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

func (cli *CelaenoClient) GetLocalUser() (*shared.User, error) {
	if cli.LocalUser == nil {
		return nil, fmt.Errorf("not logged in")
	}
	return cli.LocalUser, nil
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
