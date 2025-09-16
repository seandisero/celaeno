package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandSetDisplayName(cfg cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments passed into Set Display Name should be:\n/set displayname <new_displayname>")
	}

	// TODO: finish this one next maybe

	return nil
}
