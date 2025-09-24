package commands

import (
	"fmt"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandLeaveChat(cfg *cliapi.CelaenoConfig, args ...string) error {
	if cfg.Client.Connection != nil {
		err := cfg.Client.Connection.Close(websocket.StatusNormalClosure, "user left the chat")
		if err != nil {
			return fmt.Errorf("failed to close websocket connection: %w", err)
		}
	}

	cfg.Client.ChatRoom = ""

	cfg.Client.Screen.CelaenoResponse("you left the chat")

	return nil
}
