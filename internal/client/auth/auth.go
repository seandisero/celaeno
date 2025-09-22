package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

const (
	tokenFileName = "token.jwt"
	tokenFileMode = 0600
)

func getTokenFilePath(username string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find user home dir: %w", err)
	}
	return filepath.Join(homeDir, ".config", "celaeno", username+"_"+tokenFileName), nil
}

func makeConfigIfNotExists(filePath string) {
	dir := filepath.Dir(filepath.Clean(filePath))
	var err error
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
		if err != nil {
			slog.Error("Error making config dir")
		}
	}
}

func writeDataToTokenFile(data []byte, username string) error {
	filePath, err := getTokenFilePath(username)
	if err != nil {
		return fmt.Errorf("could not get token file path")
	}
	file, err := os.OpenFile(filepath.Clean(filePath), os.O_CREATE|os.O_TRUNC, tokenFileMode)
	if err != nil {
		slog.Error("error creating file", "error", err)
		return err
	}
	defer file.Close()

	err = os.WriteFile(filepath.Clean(filePath), data, tokenFileMode)
	if err != nil {
		return fmt.Errorf("error when writing token to file: %w", err)
	}
	return nil
}

func AuthToken(username string) (string, error) {
	return findAuthToken(username)
}

func SetAuthToken(token string, username string) error {
	return saveTokenToFile(token, username)
}

func saveTokenToFile(token string, username string) error {
	filePath, err := getTokenFilePath(username)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	makeConfigIfNotExists(filePath)
	err = writeDataToTokenFile([]byte(token), username)
	if err != nil {
		return err
	}

	return nil
}

func ApplyBearerToken(r *http.Request, username string) error {
	authToken, err := findAuthToken(username)
	if err != nil {
		return fmt.Errorf("could not get auth token: %w", err)
	}

	bearerToken := "Bearer " + authToken
	r.Header.Set("Authorization", bearerToken)
	return nil
}

func ApplyBearerTokenToHeader(header *http.Header, username string) error {
	authToken, err := findAuthToken(username)
	if err != nil {
		return fmt.Errorf("could not get auth token: %w", err)
	}

	bearerToken := "Bearer " + authToken
	header.Set("Authorization", bearerToken)
	return nil
}

func findAuthToken(username string) (string, error) {
	filePath, err := getTokenFilePath(username)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("token file does not exist")
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
