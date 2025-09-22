package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/seandisero/celaeno/internal/client/cliapi"
)

func CommandDeleteUser(cfg *cliapi.CelaenoConfig, args ...string) error {

	fmt.Printf(" > enter your password to delete user: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	password := scanner.Text()
	fmt.Println(password)

	err := cfg.Client.DeleteUser(password)
	if err != nil {
		return err
	}

	return nil
}
