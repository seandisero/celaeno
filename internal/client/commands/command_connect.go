package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seandisero/celaeno/internal/client/auth"
	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandConnect(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("need an argument to connect: /connect <username>")
	}

	cfg.Client.Connect(args[0])

	type connectionRequest struct {
		Username string `json:"username"`
	}

	data := connectionRequest{
		Username: args[0],
	}
	jsonData, err := json.Marshal(data)

	dataBuffer := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("GET", cfg.Client.URL+"/api/connect", dataBuffer)
	if err != nil {
		return fmt.Errorf("could not make new request %w", err)
	}

	auth.ApplyBearerToken(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfg.Client.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("something went wrong establishing connection")
	}
	return nil
}
