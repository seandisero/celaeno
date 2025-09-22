package cliapi

import (
	"fmt"
	"testing"
	"time"

	"github.com/seandisero/celaeno/internal/shared"
)

func TestEncryption(t *testing.T) {
	client := NewClient(time.Minute)
	client.LocalUser = &shared.User{
		Username: "sean",
		ID:       []byte("3b1e05b3-6cc4-463f-aa0b-00de16300fd3"),
	}

	message := "this is my message"

	encrepted, err := client.Encrypt([]byte(message))
	if err != nil {
		t.FailNow()
	}

	n := string(encrepted)
	nn := []byte(n)

	decrypted, err := client.Decrypt(nn)
	if err != nil {
		fmt.Printf("%v", err)
		t.FailNow()
	}

	if string(decrypted) != message {
		fmt.Println("message does not match")
		t.FailNow()
	}

	fmt.Printf("origional message: %s\n", message)
	fmt.Printf("decrypted message: %s\n", string(decrypted))
}
