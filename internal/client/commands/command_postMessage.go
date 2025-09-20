package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/shared"
)

func CommandPostMessage(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("need more arguments to post message")
	}

	if cfg.Client.LocalUser == nil {
		cfg.Client.Screen.HandleMessage(shared.Message{
			Username: "echo",
			Message:  args[0],
		})
		return fmt.Errorf("not logged in")
	}

	name := cfg.Client.LocalUser.Username
	to := cfg.Client.ChatRoom

	message := shared.Message{
		Message:  args[0],
		Username: name,
		To:       to,
	}

	err := cfg.Client.PostMessage(&message)
	if err != nil {
		return err
	}

	return nil
}
