package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yagnikpt/webhook-relay/internal/config"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func InitAuth(BASE_URL string) error {
	accessToken, err := LoginWithDeviceFlow()
	if err != nil {
		return err
	}

	err = SaveTokenToKeyring(accessToken)
	if err != nil {
		return err
	}
	return nil
}

// LoginWithDeviceFlow initiates the OAuth2 device flow and returns the access token.
func LoginWithDeviceFlow() (string, error) {
	cfg := config.Load()
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	deviceCodeURL := "https://github.com/login/device/code"
	accessTokenURL := "https://github.com/login/oauth/access_token"
	userInfoURL := "https://api.github.com/user"

	deviceCodePayload := map[string]string{
		"client_id": cfg.GitHubClientID,
		"scope":     "read:user",
	}
	body, _ := json.Marshal(deviceCodePayload)
	req, err := http.NewRequest("POST", deviceCodeURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	var deviceCodeRes DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCodeRes); err != nil {
		panic(err)
	}

	fmt.Printf("Please visit: %s\nenter the code: %s\nWaiting for the token...\n", deviceCodeRes.VerificationUri, deviceCodeRes.UserCode)

	ticker := time.NewTicker(time.Duration(deviceCodeRes.Interval)*time.Second + time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		accessTokenPayload := map[string]string{
			"client_id":   cfg.GitHubClientID,
			"device_code": deviceCodeRes.DeviceCode,
			"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		}
		body, _ := json.Marshal(accessTokenPayload)
		req, err := http.NewRequest("POST", accessTokenURL, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return "", err
		}
		defer resp.Body.Close()

		var accessTokenRes AccessTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&accessTokenRes); err != nil {
			fmt.Println("Error decoding response:", err)
			panic(err)
		}
		if accessTokenRes.AccessToken != "" {
			req, err := http.NewRequest("GET", userInfoURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return "", err
			}
			req.Header.Set("Authorization", "Bearer "+accessTokenRes.AccessToken)
			req.Header.Set("Accept", "application/vnd.github+json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				return "", err
			}
			defer resp.Body.Close()

			var userInfo map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
				return "", err
			}

			fmt.Printf("Authentication successful! Welcome, %s\n", userInfo["login"])
			return accessTokenRes.AccessToken, nil
		}
	}
}
