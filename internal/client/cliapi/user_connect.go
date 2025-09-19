package cliapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seandisero/celaeno/internal/client/auth"
)

func (cli *CelaenoClient) Connect(name string) error {
	type connectionRequest struct {
		Username string `json:"username"`
	}

	data := connectionRequest{
		Username: name,
	}
	jsonData, err := json.Marshal(data)

	dataBuffer := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", cli.URL+"/api/connect", dataBuffer)
	if err != nil {
		return fmt.Errorf("could not make new request %w", err)
	}

	auth.ApplyBearerToken(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("something went wrong establishing connection")
	}
	return nil
}
