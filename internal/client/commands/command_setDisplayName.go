package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandSetUserAttr(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments passed into Set Display Name should be:\n/set displayname <new_displayname>")
	}

	switch args[0] {
	case "displayname":
		user, err := cfg.Client.SetDisplayName(args[1])
		if err != nil {
			return err
		}

		cfg.Client.LocalUser = &user

		fmt.Println(" + ")
		fmt.Printf(" > name set to %s\n", user.Displayname.String)
		fmt.Println(" + ")
	case "cipher":
		err := cfg.Client.SetUserCipher(args[1])
		if err != nil {
			return err
		}

		fmt.Println(" + ")
		fmt.Println(" > cipher has been set")
		fmt.Println(" + ")
	}

	// TODO: finish this one next maybe

	return nil
}
