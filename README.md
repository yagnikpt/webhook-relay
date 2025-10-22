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

3. Build the binaries:
   ```bash
   make all
   ```
   This builds both the CLI client (`whrelay`) and the server (`server/whrelay_server`).

## Usage

1. Start the server:
   ```bash
   make server
   ```
   The server will run on `http://localhost:8080`.

2. In another terminal, start the client:
   ```bash
   make cli
   ```
   The client will authenticate with GitHub using device flow, save the access token, create a webhook endpoint, connect via WebSocket, and relay incoming webhooks to `http://localhost:3000` (adjust in code if needed).

3. Send webhooks:
   - External services can POST to `http://your-server:8080/webhook/{id}` with `Authorization: Bearer <github_token>` to relay webhooks to the connected client.

## API Endpoints

- `GET /`: Health check endpoint.
- `GET /webhook`: Create a new webhook endpoint and return its ID.
- `POST /webhook/{id}`: Receive and forward webhook payloads to the connected client.
- `GET /connect/{id}`: Establish a WebSocket connection for receiving relayed webhooks.

## Development

- Download dependencies: `make download`
- Tidy modules: `make tidy`
- Format code: `make fmt`
- Vet code: `make vet`
- Run tests: `make test`
- Build CLI: `make build-cli`
- Build server: `make build-server`
- Clean binaries: `make clean`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License