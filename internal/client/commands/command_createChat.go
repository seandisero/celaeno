package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/shared"
)

func CommandCreateChat(cfg *cliapi.CelaenoConfig, args ...string) error {
	err := cfg.Client.CreateChat()
	if err != nil {
		return fmt.Errorf("error creating chat: %w", err)
	}

	message := shared.Message{
		Username: "celaeno",
		To:       cfg.Client.LocalUser.Username,
		Message:  "created chat, other users can connect to you using the connect command",
		Incoming: true,
	}
	cfg.Client.Screen.HandleMessage(message)

	return nil
}
