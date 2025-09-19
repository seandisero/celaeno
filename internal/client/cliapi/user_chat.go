package cliapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seandisero/celaeno/internal/shared"
)

func (cli *CelaenoClient) PostMessage(message *shared.Message) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal data")
	}

	data := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", cli.URL+"/api/chat/publish", data)
	if err != nil {
		return fmt.Errorf("unsucessful creation of request: %w", err)
	}

	resp, err := cli.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unsucessful request: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("error code: %d", resp.StatusCode)
	}

	return nil
}
