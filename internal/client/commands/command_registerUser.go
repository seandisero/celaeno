package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandRegisterUser(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough args to register user should be \nregister <name> <password>")
	}

	name := args[0]
	password := args[1]

	user, err := cfg.Client.RegisterUser(name, password)
	if err != nil {
		return fmt.Errorf("could not register user: %w", err)
	}

	fmt.Println(" + ")
	fmt.Printf(" > rigistered new user: %s\n", user.Username)
	fmt.Println(" > time to log in: us /login <username> <password>")
	fmt.Println(" + ")

	return nil
}
