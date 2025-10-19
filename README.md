# Webhook Relay

A lightweight, self-hosted webhook relay server written in Go. Similar to ngrok or smee.io, it allows external services (like GitHub, Stripe, etc.) to send webhooks to your local development server securely.

## Features

- **Webhook Forwarding**: Receive webhooks from external services and forward them to connected clients via WebSocket.
- **Real-time Delivery**: Maintain persistent WebSocket connections for instant webhook delivery.
- **Simple API**: Easy-to-use endpoints for creating endpoints and connecting clients.
- **Thread-safe**: Uses `sync.Map` for concurrent connection management.
- **Lightweight**: Minimal dependencies, built with Go's standard library and `github.com/coder/websocket`.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yagnikpt/webhook-relay.git
   cd webhook-relay
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the server:
   ```bash
   go build -o whrelay ./cmd/whrelay
   ```

## Usage

1. Start the server:
   ```bash
   ./whrelay
   ```
   The server will run on `http://localhost:8080`.

2. Create a webhook endpoint:
   - Visit `http://localhost:8080/webhook` to generate a unique ID.

3. Connect a client:
   - Use a WebSocket client to connect to `ws://localhost:8080/connect/{id}`.

4. Send webhooks:
   - External services can POST to `http://your-server:8080/webhook/{id}` to relay webhooks to the connected client.

## API Endpoints

- `GET /`: Health check endpoint.
- `GET /webhook`: Create a new webhook endpoint and return its ID.
- `POST /webhook/{id}`: Receive and forward webhook payloads to the connected client.
- `GET /connect/{id}`: Establish a WebSocket connection for receiving relayed webhooks.

## Development

- Run tests: `go test ./...`
- Build: `make build` (if Makefile is configured)

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License