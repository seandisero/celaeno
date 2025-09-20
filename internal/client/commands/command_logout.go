package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/shared"
)

func CommandLogout(cfg *cliapi.CelaenoConfig, args ...string) error {

	// TODO: eventually I'll want to patition the server to deatroy any web sockets connected to my pc
	// but for now I'll just remove the data form the token.jwt in the config file.

	err := cfg.Client.Logout()
	if err != nil {
		return fmt.Errorf("could not log out: %w", err)
	}

	cfg.Client.Cancel()

	message := shared.Message{
		Username: "celaeno",
		Incoming: true,
		Message:  "you are logged out",
	}

	cfg.Client.Screen.HandleMessage(message)

	return nil
}
