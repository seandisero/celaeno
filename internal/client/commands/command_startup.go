package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandStartup(cfg *cliapi.CelaenoConfig, args ...string) error {
	cfg.Client.Screen.CelaenoResponse("I'll try and start up the server for you")
	go timedStartup(cfg)
	return nil
}

func timedStartup(cfg *cliapi.CelaenoConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.Client.URL+"/startup", http.NoBody)
	if err != nil {
		slog.Error("error creating new request with context", "error", err)
		return
	}
	resp, err := cfg.Client.HttpClient.Do(req)
	if err != nil {
		slog.Error("error doing http request", "error", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading responce data from server")
		return
	}

	type response struct {
		Status string `json:"status"`
	}
	var jsonResp response
	err = json.Unmarshal(data, &jsonResp)
	if err != nil {
		slog.Error("coul not unmarshal data from response", "error", err)
		return
	}

	cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("looks like the server is: %s", jsonResp.Status))
}
