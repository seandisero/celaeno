package commands

import (
	"encoding/base64"
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

	if cfg.Client.ChatRoom == "" {
		return fmt.Errorf("not a part of any chat room")
	}

	encryptedMessage, err := cfg.Client.Encrypt([]byte(args[0]))
	if err != nil {
		return err
	}

	name := cfg.Client.LocalUser.Username
	if cfg.Client.LocalUser.Displayname.Valid {
		name = cfg.Client.LocalUser.Displayname.String
	}
	to := cfg.Client.ChatRoom

	encrypted := base64.StdEncoding.EncodeToString(encryptedMessage)

	message := shared.Message{
		Message:  encrypted,
		Username: name,
		To:       to,
	}

	err = cfg.Client.PostMessage(&message)
	if err != nil {
		return err
	}

	return nil
}
