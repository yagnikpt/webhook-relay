package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

func SaveTokenToKeyring(token string) error {
	err := keyring.Set("webhook_relay", "access_token", token)
	if err != nil {
		return fmt.Errorf("failed to save token to keyring: %w", err)
	}
	return nil
}

func GetTokenFromKeyring() (string, error) {
	secret, err := keyring.Get("webhook_relay", "access_token")
	if err != nil {
		return "", fmt.Errorf("failed to get token from keyring: %w", err)
	}
	return secret, nil
}

func DeleteTokenFromKeyring() error {
	err := keyring.Delete("webhook_relay", "access_token")
	if err != nil {
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}
