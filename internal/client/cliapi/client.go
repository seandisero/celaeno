package cliapi

import (
	"net/http"
	"time"
)

type CelaenoClient struct {
	HttpClient *http.Client
	URL        string
}

func NewClient(timeout time.Duration) CelaenoClient {
	client := CelaenoClient{
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return client
}
