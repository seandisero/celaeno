package commands

import (
	"fmt"
	"strings"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandListChatRooms(cfg *cliapi.CelaenoConfig, args ...string) error {
	chatRooms, err := cfg.Client.ListChatRooms()
	if err != nil {
		return err
	}

	cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("available chat rooms are: %s", strings.Join(chatRooms, ", ")))
	return nil
}
