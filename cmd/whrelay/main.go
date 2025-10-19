package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/coder/websocket"
)

var BASE_URL string
var WS_PROTOCOL string
var HTTP_PROTOCOL string
var PORT string
var FORWARD_ENDPOINT string

func getEndpoint() string {
	url := fmt.Sprintf("%s%s/webhook", HTTP_PROTOCOL, BASE_URL)
	resp, err := http.Get(url)
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

	conn, _, err := websocket.Dial(ctx, endpoint, nil)
	if err != nil {
		fmt.Printf("%d", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	receiveEndpoint := fmt.Sprintf("%s%s/webhook/%s", HTTP_PROTOCOL, BASE_URL, id)
	fmt.Println("Listening for webhook on endpoint: " + receiveEndpoint + "\n")

	for {
		msgtype, msg, err := conn.Read(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("Received webhook payload: %s\n", string(msg))
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

			localUrl := "http://localhost:" + PORT + FORWARD_ENDPOINT

			req, err := http.NewRequest("POST", localUrl, bytes.NewBufferString(body))
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
				fmt.Println("Error sending request to local server:", err)
				continue
			}
			resp.Body.Close()

			fmt.Printf("Sent webhook payload to %s\n", localUrl)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: whrelay <local-port> <endpoint>")
		return
	}
	PORT = os.Args[1]
	FORWARD_ENDPOINT = os.Args[2]
	dev := os.Getenv("ENVIRONMENT") == "development"
	if dev {
		BASE_URL = "localhost:8080"
		WS_PROTOCOL = "ws://"
		HTTP_PROTOCOL = "http://"
	} else {
		BASE_URL = "whrelay.example.com"
		WS_PROTOCOL = "wss://"
		HTTP_PROTOCOL = "https://"
	}

	endpointID := getEndpoint()
	connectWebSocket(endpointID)
}
