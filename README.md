# Webhook Relay

A lightweight, self-hosted webhook relay server written in Go. Similar to ngrok or smee.io, it allows external services (like GitHub, Stripe, etc.) to send webhooks to your local development server securely.

https://github.com/user-attachments/assets/971e20fd-c11a-4efb-a876-35771962051d

## Features

- **Webhook Forwarding**: Receive webhooks from external services and forward them to connected clients via WebSocket.
- **Real-time Delivery**: Maintain persistent WebSocket connections for instant webhook delivery.
- **Simple API**: Easy-to-use endpoints for creating endpoints and connecting clients.
- **Thread-safe**: Uses `sync.Map` for concurrent connection management.
- **Lightweight**: Minimal dependencies, built with Go's standard library and `github.com/coder/websocket`.

## Installation

### Using Go
```bash
go install github.com/yagnikpt/webhook-relay/cmd/whrelay@latest
```

### From the source

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
   make build-cli
   ```
   This only builds the CLI client (`whrelay`)

## Usage

1. Login (github oauth) by running:
   ```bash
   whrelay login
   ```
   - The client will authenticate with GitHub using device flow, save the access token.
2. Start listening webhooks:
   ```bash
   whrelay <local-port> <local-endpoint>
   ```
   - This creates a webhook endpoint, connect via WebSocket, and relay incoming webhooks to localhost:`<local-port>`/`<local-endpoint>`
3. Send webhooks:
   - External services can POST to `https://wh-relay.azurewebsites.net/webhook/{id}`

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
