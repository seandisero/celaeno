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

func getTokenFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find user home dir: %w", err)
	}
	return filepath.Join(homeDir, ".config", "celaeno", tokenFileName), nil
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

func writeDataToTokenFile(data []byte) error {
	filePath, err := getTokenFilePath()
	if err != nil {
		return fmt.Errorf("could not get token file path")
	}
	file, err := os.OpenFile(filepath.Clean(filePath), os.O_CREATE|os.O_TRUNC, tokenFileMode)
	if err != nil {
		slog.Info("error creating file")
		return err
	}
	defer file.Close()

	err = os.WriteFile(filepath.Clean(filePath), data, tokenFileMode)
	if err != nil {
		return fmt.Errorf("error when writing token to file: %w", err)
	}
	return nil
}

func AuthToken() (string, error) {
	return findAuthToken()
}

func SetAuthToken(token string) error {
	return saveTokenToFile(token)
}

func saveTokenToFile(token string) error {
	filePath, err := getTokenFilePath()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	makeConfigIfNotExists(filePath)
	err = writeDataToTokenFile([]byte(token))

	return nil
}

func ApplyBearerToken(r *http.Request) error {
	authToken, err := findAuthToken()
	if err != nil {
		return fmt.Errorf("could not get auth token: %w", err)
	}

	bearerToken := "Bearer " + authToken
	r.Header.Set("Authorization", bearerToken)
	return nil
}

func findAuthToken() (string, error) {
	filePath, err := getTokenFilePath()
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
