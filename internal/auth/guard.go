package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func GetUserName() (string, error) {
	secret, err := keyring.Get("webhook_relay", "access_token")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println("Error creating request user info", "error", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error validating token", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Invalid GitHub token")
	}

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", fmt.Errorf("Failed to parse user info: %w", err)
	}
	return userInfo["login"].(string), nil
}
