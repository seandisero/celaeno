package cliapi

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coder/websocket"
)

type CelaenoConfig struct {
	Client   *CelaenoClient
	Commands map[string]func(cfg *CelaenoConfig, args ...string) error
}

var LINE_DELIMINATOR = ">"

func getCommandString(s string) (string, []string, error) {
	input := strings.Split(s, " ")
	if len(input) < 1 {
		return "", make([]string, 0), fmt.Errorf("could not split input string for command")
	}
	if len(input) == 1 {
		return input[0], make([]string, 0), nil
	}

	return input[0], input[1:], nil
}

func (cfg *CelaenoConfig) ExitApplication(exitSignal chan os.Signal) {
	<-exitSignal

	if cfg.Client.Connection != nil {
		fmt.Println("closing client connection")
		err := cfg.Client.Connection.Close(websocket.StatusNormalClosure, "user is closing the program")
		if err != nil {
			slog.Error("error occerred in connection", "error", err)
		}
	}

	fmt.Printf(" > g")
	for _, c := range "oodbuy" {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%s", string(c))
	}
	fmt.Println()
	os.Exit(0)
}

func (cfg *CelaenoConfig) doCommandPostMessage(message string) {
	postCommand, ok := cfg.Commands["post-message"]
	if !ok {
		cfg.Client.Screen.CelaenoResponse("internal error getting post-message command")
		return
	}
	err := postCommand(cfg, message)
	if err != nil {
		if strings.Contains(err.Error(), "no authorization token") {
			cfg.Client.Screen.CelaenoResponse("it looks like your auth token doesn't work any more, try loggin in again")
			cfg.Client.Screen.CelaenoResponse("to log in use /login <username> <password>")
			return
		} else if strings.Contains(err.Error(), "expired") {
			cfg.Client.Screen.CelaenoResponse("login timed out, please login again")
			cfg.Client.Screen.CelaenoResponse("to log in use /login <username> <password>")
			return
		} else if strings.Contains(err.Error(), "not logged in") {
			cfg.Client.Screen.CelaenoResponse("not logged in")
			cfg.Client.Screen.CelaenoResponse("to log in use /login <username> <password>")
			return
		} else if strings.Contains(err.Error(), "you must make a connection or") {
			cfg.Client.Screen.CelaenoResponse("make a connection or start a chat")
			cfg.Client.Screen.CelaenoResponse("/create-chat")
			cfg.Client.Screen.CelaenoResponse("/connect <username>")
			return
		} else if strings.Contains(err.Error(), "chat no longer exists") {
			cfg.Client.Screen.CelaenoResponse("chat no longer exists, the host probably left.")
		}
		cfg.Client.Screen.CelaenoResponse("some error occerred")
		cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("some error occerred: %v", err))
	}
}

func (cfg *CelaenoConfig) doCommand(message string) {
	commandStr, args, err := getCommandString(strings.Trim(message, "/"))
	if err != nil {
		cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("error parsing the command: %s", commandStr))
		return
	}
	command, ok := cfg.Commands[commandStr]
	if !ok {
		cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("the command %s does not exist", commandStr))
		return
	}
	err = command(cfg, args...)
	if err != nil {
		cfg.Client.Screen.CelaenoResponse(fmt.Sprintf("error running command: %s: %v", commandStr, err.Error()))
		return
	}
}

func StartRepl(cfg *CelaenoConfig) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go cfg.ExitApplication(sigs)

	scanner := bufio.NewScanner(os.Stdin)
	cfg.Client.Screen.PrintWelcomeMessage()
	for range cfg.Client.Screen.Width {
		fmt.Printf("-")
	}
	fmt.Println()

	for {
		if !scanner.Scan() {
			break
		}
		message := scanner.Text()

		err := cfg.Client.Screen.ClearInput(message)
		if err != nil {
			slog.Error("error clearing input", "error", err)
			return
		}

		if message == "" {
			continue
		}

		if message[0] == '/' {
			cfg.doCommand(message)
			continue
		}
		cfg.doCommandPostMessage(message)
	}
}
