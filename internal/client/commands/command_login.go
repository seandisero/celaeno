package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandLogin(cfg cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough args for login: example \n/login <name> <password>")
	}

	name := args[0]
	password := args[1]

	user, err := cfg.Client.Login(name, password)
	if err != nil {
		return fmt.Errorf("error loggin in: %w", err)
	}

	fmt.Printf(" > logged in as %s\n", user.Username)

	return nil
}
