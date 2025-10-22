package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/yagnikpt/webhook-relay/internal/auth"
	"github.com/yagnikpt/webhook-relay/internal/config"
	"github.com/yagnikpt/webhook-relay/internal/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	godotenv.Load()
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", server.HelloHandler)
	mux.Handle("GET /webhook", auth.Middleware(http.HandlerFunc(server.CreateWebhookHandler)))
	mux.Handle("POST /webhooks/{id}", auth.Middleware(http.HandlerFunc(server.ReceiveWebhookHandler)))
	mux.Handle("GET /connect/{id}", auth.Middleware(http.HandlerFunc(server.ConnectHandler)))
	http.ListenAndServe(":"+cfg.Port, mux)
}
