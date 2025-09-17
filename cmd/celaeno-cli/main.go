package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/client/commands"
)

func mapCommands(cfg cliapi.CelaenoConfig) {
	cfg.Commands["post-message"] = commands.CommandPostMessage
	cfg.Commands["login"] = commands.CommandLogin
	cfg.Commands["logout"] = commands.CommandLogout
	cfg.Commands["whoami"] = commands.CommandGetUser
	cfg.Commands["register"] = commands.CommandRegisterUser
	cfg.Commands["set displayname"] = commands.CommandSetDisplayName
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Could not open .env", "ERROR", err)
	}

	celaenoClient := cliapi.NewClient(5 * time.Second)
	celaenoClient.URL = os.Getenv("SERVER_URL") + os.Getenv("PORT")

	celaenoConfig := cliapi.CelaenoConfig{
		Client:   celaenoClient,
		Commands: make(map[string]func(cfg cliapi.CelaenoConfig, args ...string) error),
	}

	mapCommands(celaenoConfig)
	cliapi.StartRepl(celaenoConfig)
}
