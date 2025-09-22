package cliapi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (cli *CelaenoClient) Encrypt(message []byte) ([]byte, error) {
	key, err := getCipherKey(cli.LocalUser.ID)
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return []byte{}, err
	}

	encryptedMessage := aesgcm.Seal(nil, nonce, message, nil)

	result := append(nonce, encryptedMessage...)

	return result, nil
}

func (cli *CelaenoClient) Decrypt(message []byte) ([]byte, error) {
	key, err := getCipherKey(cli.LocalUser.ID)
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(message) < nonceSize {
		return []byte{}, fmt.Errorf("message too short")
	}

	nonce := message[:nonceSize]
	encryptedMessage := message[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		// slog.Error("error decrypting message", "error", err)
		return encryptedMessage, err
	}

	return decrypted, nil
}

func (cli *CelaenoClient) SetUserCipher(cipherKey string) error {
	path, err := getCipherKeyFilePath(cli.LocalUser.ID)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	hasher.Write([]byte(cipherKey))
	hashBytes := hasher.Sum(nil)

	err = os.WriteFile(path, hashBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (cli *CelaenoClient) SetupCipher() error {
	var err error
	block, err := newCipherBlock(cli.LocalUser.ID)
	if err != nil {
		return err
	}
	cli.Block = &block
	return nil
}

func newCipherBlock(userID []byte) (cipher.Block, error) {
	key, err := getCipherKey(userID)
	if err != nil {
		return nil, err
	}
	if key == nil {
		fmt.Println(" > no cipher key set: use /set cipher <key>")
		return nil, nil
	}
	return aes.NewCipher(key)
}

func getCipherKey(userID []byte) ([]byte, error) {
	path, err := getCipherKeyFilePath(userID)
	if err != nil {
		return []byte{}, err
	}

	// file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_TRUNC, 0600)
	// if err != nil {
	// return []byte{}, err
	// }
	// file.Close()

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func getCipherKeyFilePath(id []byte) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find user home dir: %w", err)
	}
	return filepath.Join(homeDir, ".config", "celaeno", string(id)+"_cipher"), nil
}
