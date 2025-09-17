package cliapi

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type CelaenoConfig struct {
	Client   CelaenoClient
	Commands map[string]func(cfg CelaenoConfig, args ...string) error
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

func StartRepl(cfg CelaenoConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(" %s ", LINE_DELIMINATOR)
		if !scanner.Scan() {
			break
		}
		message := scanner.Text()

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
				fmt.Println(" > ")
				fmt.Println(" > you must be logged in")
				fmt.Println(" > ")
				continue
			} else if strings.Contains(err.Error(), "needs more arguments") {
				fmt.Println(" > ")
			} else if strings.Contains(err.Error(), "expired") {
				fmt.Printf(" + \n > login timed out\n + \n")
				continue
			}
			slog.Error("posting message", "error", err)
			return
		}
	}
}
