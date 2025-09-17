package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandLogout(cfg cliapi.CelaenoConfig, args ...string) error {

	// TODO: eventually I'll want to patition the server to deatroy any web sockets connected to my pc
	// but for now I'll just remove the data form the token.jwt in the config file.

	err := cfg.Client.Logout()
	if err != nil {
		return fmt.Errorf("could not log out: %w", err)
	}

	fmt.Println(" > ")
	fmt.Println(" > you are now logged out")
	fmt.Println(" > ")

	return nil
}
