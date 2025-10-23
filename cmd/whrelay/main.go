package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/joho/godotenv"
	"github.com/yagnikpt/webhook-relay/internal/auth"
	"github.com/yagnikpt/webhook-relay/internal/utils"
)

var BASE_URL string
var WS_PROTOCOL string
var HTTP_PROTOCOL string
var PORT string = "3000"
var FORWARD_ENDPOINT string = "/"

func getEndpoint() string {
	token, err := auth.GetTokenFromKeyring()
	if err != nil {
		fmt.Println("Error retrieving token from keyring:", err)
		return ""
	}
	url := fmt.Sprintf("%s%s/webhook", HTTP_PROTOCOL, BASE_URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v", err)
			return ""
		}
		bodyString := string(bodyBytes)
		return bodyString
	} else {
		fmt.Println("Failed to hit create endpoint")
	}
	return ""
}

func connectWebSocket(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	endpoint := fmt.Sprintf("%s%s/connect/%s", WS_PROTOCOL, BASE_URL, id)

	headers := http.Header{}
	token, err := auth.GetTokenFromKeyring()
	if err != nil {
		fmt.Println("Error retrieving token from keyring:", err)
		return
	}
	headers.Add("Authorization", "Bearer "+token)

	opts := &websocket.DialOptions{
		HTTPHeader: headers,
	}

	conn, _, err := websocket.Dial(ctx, endpoint, opts)
	if err != nil {
		fmt.Printf("%d", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	userName, err := auth.GetUserName()
	if err != nil {
		fmt.Println("Error retrieving username:", err)
		return
	}

	receiveEndpoint := fmt.Sprintf("%s%s/webhook/%s", HTTP_PROTOCOL, BASE_URL, id)
	forwardEndpoint := "http://localhost:" + PORT + FORWARD_ENDPOINT

	utils.PrintInitialTUI(userName, receiveEndpoint, forwardEndpoint)

	for {
		msgtype, msg, err := conn.Read(context.Background())
		if err != nil {
			break
		}
		if msgtype == websocket.MessageText {
			var data map[string]any
			err := json.Unmarshal(msg, &data)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				continue
			}
			headers, ok := data["headers"].(map[string]any)
			if !ok {
				fmt.Println("Error parsing headers")
				continue
			}
			body, ok := data["body"].(string)
			if !ok {
				fmt.Println("Error parsing body")
				continue
			}

			req, err := http.NewRequest("POST", forwardEndpoint, bytes.NewBufferString(body))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}

			for key, value := range headers {
				if strValue, ok := value.(string); ok {
					req.Header.Set(key, strValue)
				}
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("ERROR", "Connection Refused", forwardEndpoint)
				continue
			}
			resp.Body.Close()
			fmt.Println("POST", FORWARD_ENDPOINT, resp.Status)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: whrelay <local-port> <local-endpoint>")
		return
	}
	godotenv.Load()
	dev := os.Getenv("ENVIRONMENT") == "development"
	if dev {
		BASE_URL = "localhost:8080"
		WS_PROTOCOL = "ws://"
		HTTP_PROTOCOL = "http://"
	} else {
		BASE_URL = "wh-relay.azurewebsites.net"
		WS_PROTOCOL = "wss://"
		HTTP_PROTOCOL = "https://"
	}

	if os.Args[1] == "login" {
		err := auth.InitAuth(HTTP_PROTOCOL + BASE_URL)
		if err != nil {
			fmt.Println("Login failed:", err)
		} else {
			fmt.Println("Usage: whrelay <local-port> <local-endpoint>")
		}
		return
	}
	if err := auth.AuthGuard(); err != nil {
		return
	}
	if len(os.Args) < 3 {
		fmt.Println("Usage: whrelay <local-port> <local-endpoint>")
		return
	}
	PORT = os.Args[1]
	_, err := strconv.Atoi(PORT)
	if err != nil {
		fmt.Println("Invalid port number")
		return
	}
	FORWARD_ENDPOINT = os.Args[2]
	if FORWARD_ENDPOINT[0] != '/' {
		fmt.Println("Endpoint should start with '/'")
		return
	}

	endpointID := getEndpoint()
	if endpointID == "" {
		return
	}
	connectWebSocket(endpointID)
}
