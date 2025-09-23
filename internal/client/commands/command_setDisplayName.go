package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandSetUserAttr(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		cfg.Client.Screen.CelaenoResponse("not enough arguments passed into set")
		cfg.Client.Screen.CelaenoResponse("/set <attribute> <new_displayname>")
		return fmt.Errorf("not enough arguments passed into set")
	}

	switch args[0] {
	case "displayname":
		user, err := cfg.Client.SetDisplayName(args[1])
		if err != nil {
			return err
		}

		cfg.Client.LocalUser = &user
		cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("display name set to: %s", user.Displayname.String))
	case "cipher":
		err := cfg.Client.SetUserCipher(args[1])
		if err != nil {
			return err
		}
		cfg.Client.Screen.CelaenoResponse("new cipher has been set")
	}

	return nil
}
