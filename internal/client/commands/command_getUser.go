package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandGetUser(cfg cliapi.CelaenoConfig, args ...string) error {
	user, err := cfg.Client.GetUser()
	if err != nil {
		return fmt.Errorf("error getting user info: %w", err)
	}

	fmt.Printf(" %s \n", cliapi.LINE_DELIMINATOR)
	fmt.Printf(" %s username: %s\n", cliapi.LINE_DELIMINATOR, user.Username)
	if user.Displayname.Valid {
		fmt.Printf(" %s display name: %s\n", cliapi.LINE_DELIMINATOR, user.Displayname.String)
	}

	return nil
}
