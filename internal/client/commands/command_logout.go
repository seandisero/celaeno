package commands

import (
	"fmt"

	"github.com/coder/websocket"
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

	if cfg.Client.Connection != nil {
		err = cfg.Client.Connection.Close(websocket.StatusNormalClosure, "user has logged out")
		if err != nil {
			return fmt.Errorf("could not close connection to websocket: %w", err)
		}
	}

	message := shared.Message{
		Username: "celaeno",
		Incoming: true,
		Message:  "you are logged out",
	}

	cfg.Client.Screen.HandleMessage(message)
	cfg.Client.ChatRoom = ""

	return nil
}
