package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/lithammer/shortuuid/v4"
)

var connections sync.Map

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func createWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	id := shortuuid.New()
	connections.Store(id, nil)
	fmt.Fprintf(w, "Webhook endpoint created with ID: %s", id)
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
	if conn, ok := connections.Load(id); ok {
		err := conn.(*websocket.Conn).Write(context.Background(), websocket.MessageText, payload)
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
	fmt.Fprintf(w, "Connecting to webhook endpoint with ID: %s", id)
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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", hello)
	mux.HandleFunc("GET /webhook", createWebhookEndpoint)
	mux.HandleFunc("POST /webhook/{id}", receiveWebhook)
	mux.HandleFunc("GET /connect/{id}", handleConnect)
	http.ListenAndServe(":8080", mux)
}
