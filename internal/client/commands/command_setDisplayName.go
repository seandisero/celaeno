package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
	"github.com/seandisero/celaeno/internal/shared"
)

func CommandSetUserAttr(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments passed into Set Display Name should be:\n/set displayname <new_displayname>")
	}

	var user shared.User
	var err error
	switch args[0] {
	case "displayname":
		user, err = cfg.Client.SetDisplayName(args[1])
	}
	if err != nil {
		return err
	}

	fmt.Println(" + ")
	fmt.Printf(" > display name set to %s", user.Displayname.String)
	fmt.Println(" + ")

	// TODO: finish this one next maybe

	return nil
}
