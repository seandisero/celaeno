package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandConnect(cfg *cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("need an argument to connect: /connect <username>")
	}

	err := cfg.Client.Connect(args[0])
	if err != nil {
		return fmt.Errorf("could not make connection with %s: %w", args[0], err)
	}

	return nil
}
