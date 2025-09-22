package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/client/commands"
	"github.com/seandisero/celaeno/internal/client/screen"
)

func mapCommands(cfg *cliapi.CelaenoConfig) {
	cfg.Commands["login"] = commands.CommandLogin
	cfg.Commands["logout"] = commands.CommandLogout
	cfg.Commands["whoami"] = commands.CommandGetUser
	cfg.Commands["register"] = commands.CommandRegisterUser
	cfg.Commands["deleteme"] = commands.CommandDeleteUser

	cfg.Commands["set"] = commands.CommandSetUserAttr

	cfg.Commands["connect"] = commands.CommandConnect
	cfg.Commands["create-chat"] = commands.CommandCreateChat
	cfg.Commands["post-message"] = commands.CommandPostMessage
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Could not open .env", "ERROR", err)
	}

	celaenoClient := cliapi.NewClient(5 * time.Second)
	url := os.Getenv("SERVER_URL")
	port := os.Getenv("PORT")
	celaenoClient.URL = "http" + url + port
	// celaenoClient.WS_URL = "ws" + url + port

	celaenoClient.Screen = screen.NewScreen(64)

	ctx, cancel := context.WithCancel(context.Background())
	celaenoClient.Screen.Cancel = cancel
	go celaenoClient.Screen.MessageLoop(ctx)

	celaenoConfig := cliapi.CelaenoConfig{
		Client:   celaenoClient,
		Commands: make(map[string]func(cfg *cliapi.CelaenoConfig, args ...string) error),
	}

	mapCommands(&celaenoConfig)
	cliapi.StartRepl(&celaenoConfig)
}
