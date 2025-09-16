package commands

import (
	"fmt"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandPostMessage(cfg cliapi.CelaenoConfig, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("need more arguments to post message")
	}

	err := cfg.Client.PostMessage(args[0])
	if err != nil {
		return err
	}

	return nil
}
