package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

func AuthGuard() error {
	secret, err := keyring.Get("webhook_relay", "access_token")
	if err != nil || secret == "" {
		fmt.Println("You are not logged in. Please run\n\n\twhrelay login\n\nto authenticate.")
		return fmt.Errorf("authentication required")
	}
	return nil
}
