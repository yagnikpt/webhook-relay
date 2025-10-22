package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/coder/websocket"
	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v4"
)

var connections sync.Map

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func createWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user_id").(string)
	dev := os.Getenv("ENVIRONMENT") == "development"
	var id string
	if dev {
		id = "test-webhook-id"
	} else {
		if _, ok := connections.Load(user_id); ok {
			id = shortuuid.New()
		} else {
			id = user_id
		}
	}
	connections.Store(id, nil)
	fmt.Fprintf(w, "%s", id)
}

func receiveWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, exists := connections.Load(id)
	if !exists {
		http.Error(w, "Webhook endpoint not found", http.StatusNotFound)
		return
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}
	message := map[string]any{
		"headers": headers,
		"body":    string(payload),
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Failed to marshal message", http.StatusInternalServerError)
		return
	}
	if conn, ok := connections.Load(id); ok {
		err := conn.(*websocket.Conn).Write(context.Background(), websocket.MessageText, jsonData)
		if err != nil {
			log.Println("WebSocket write error:", err)
			http.Error(w, "Failed to send to client", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Webhook forwarded to client %s", id)
	} else {
		http.Error(w, "Client not connected", http.StatusNotFound)
	}
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, exists := connections.Load(id)
	if !exists {
		http.Error(w, "Webhook endpoint not found", http.StatusNotFound)
		return
	}
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("WebSocket accept error:", err)
		return
	}
	connections.Store(id, conn)
	defer func() {
		connections.Delete(id)
		conn.Close(websocket.StatusNormalClosure, "")
	}()
	for {
		_, _, err := conn.Read(context.Background())
		if err != nil {
			break
		}
	}
}

func main() {
	godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", hello)
	mux.Handle("GET /webhook", authMiddleware(http.HandlerFunc(createWebhookEndpoint)))
	mux.HandleFunc("POST /webhook/{id}", receiveWebhook)
	mux.Handle("GET /connect/{id}", authMiddleware(http.HandlerFunc(handleConnect)))
	http.ListenAndServe(":8080", mux)
}
