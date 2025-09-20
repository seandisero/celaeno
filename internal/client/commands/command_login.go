package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/shared"
)

func CommandLogin(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough args for login: example \n/login <name> <password>")
	}

	name := args[0]
	password := args[1]

	user, err := cfg.Client.Login(name, password)
	if err != nil {
		return fmt.Errorf("error loggin in: %w", err)
	}

	username := user.Username
	if user.Displayname.Valid {
		username = user.Displayname.String
	}

	message := shared.Message{
		Username: "celaeno",
		To:       cfg.Client.LocalUser.Username,
		Message:  fmt.Sprintf("you are logged in as %s", username),
		Incoming: true,
	}

	cfg.Client.Screen.HandleMessage(message)

	return nil
}
