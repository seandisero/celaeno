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
	"github.com/seandisero/celaeno/internal/client/screen"
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

	cfg.Client.Screen.Cancel()
	if cfg.Client.Connection != nil {
		fmt.Println("closing client connection")
		cfg.Client.Connection.Close(websocket.StatusNormalClosure, "user is closing the program")
	}

	fmt.Printf(" > g")
	for _, c := range "oodbuy" {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("%s", string(c))
	}
	fmt.Println()
	os.Exit(0)
}

func StartRepl(cfg *CelaenoConfig) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go cfg.ExitApplication(sigs)

	scanner := bufio.NewScanner(os.Stdin)

	for range cfg.Client.Screen.Width {
		fmt.Printf("-")
	}
	fmt.Println()

	for {

		if !scanner.Scan() {
			fmt.Println("breaking")
			break
		}
		message := scanner.Text()

		screen.ClearInput(message)

		if message == "" {
			continue
		}

		if message[0] == '/' {
			commandStr, args, err := getCommandString(strings.Trim(message, "/"))
			if err != nil {
				fmt.Printf(" > error parsing command: %v\n", err)
				continue
			}
			command, ok := cfg.Commands[commandStr]
			if !ok {
				fmt.Printf(" > that command does not exist: %s\n", commandStr)
				continue
			}
			err = command(cfg, args...)
			if err != nil {
				fmt.Printf(" > error running command %s: %v\n", commandStr, err)
				continue
			}
			continue
		}

		postCommand, ok := cfg.Commands["post-message"]
		if !ok {
			fmt.Println(" > could not get command Post Message")
			continue
		}
		err := postCommand(cfg, message)
		if err != nil {
			if strings.Contains(err.Error(), "no authorization token") {
				fmt.Println(" > you must be logged in")
				continue
			} else if strings.Contains(err.Error(), "needs more arguments") {
				fmt.Println(" > ")
			} else if strings.Contains(err.Error(), "expired") {
				fmt.Printf(" + \n > login timed out\n + \n")
				continue
			} else if strings.Contains(err.Error(), "not logged in") {
				continue
			}
			slog.Error("posting message", "error", err)
			continue
		}
	}
}
