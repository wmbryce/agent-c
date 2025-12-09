# Agent-C

An AI model aggregation and blockchain integration platform built with Go and Fiber.

<img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;<img src="https://img.shields.io/badge/license-Apache_2.0-red?style=for-the-badge&logo=none" alt="license" />

Agent-C provides a unified API for interacting with multiple AI model providers while managing access through blockchain-based authentication and payment systems.

## Features

- **Unified AI Model API** - Single interface for multiple AI providers (OpenAI, and more)
- **Blockchain Integration** - Ethereum smart contract support with Go bindings
- **Token-based Access Control** - API key management with token accounting
- **Multi-provider Support** - Extensible architecture for adding new AI providers
- **Smart Contract Ready** - ERC20 and custom contract support
- **Production Ready** - Built with Fiber framework for high performance
- **Auto-generated API Docs** - Swagger/OpenAPI documentation

## Tech Stack

- **Framework**: [Fiber v2](https://gofiber.io/) - Express-inspired web framework
- **Database**: PostgreSQL with pgx/v5
- **Cache**: Redis
- **Blockchain**: Ethereum via go-ethereum
- **AI**: OpenAI API integration
- **Auth**: JWT tokens
- **Migrations**: Goose
- **Logging**: zerolog
- **Validation**: go-playground/validator

## Quick Start

### Prerequisites

- Go 1.25.4+
- Docker & Docker Compose

### 1. Clone and Configure

```bash
git clone https://github.com/wmbryce/agent-c.git
cd agent-c
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start with Docker Compose (Recommended)

```bash
# Start all services (PostgreSQL, Redis, Backend)
make compose.up

# View logs
make compose.logs

# Stop services
make compose.down
```

### 3. Access the API

- **API**: http://localhost:8080
- **Swagger Docs**: http://localhost:8080/swagger/index.html

## Development

### Hot Reload Development

```bash
# Install Air for hot reloading
make air.install

# Start with hot reload
make dev
```

### Database Migrations

```bash
# Apply migrations
make goose.up

# Rollback last migration
make goose.down

# Check migration status
make goose.status

# Create new migration
make goose.create name=migration_name
```

### Manual Setup (without Docker Compose)

```bash
# Start PostgreSQL and Redis
make docker.postgres
make docker.redis

# Run migrations
make goose.up

# Generate Swagger docs
make swag

# Run application
go run cmd/main.go
```

## Project Structure

```
agent-c/
├── cmd/                      # Application entry point
│   ├── main.go              # Main entry
│   └── configs/             # Server configs
├── app/                      # Business logic
│   ├── middleware/          # HTTP middleware
│   ├── types/               # Domain models
│   ├── service/             # Business logic
│   │   └── model-providers/ # AI provider integrations
│   ├── routes/              # Route definitions
│   ├── store/               # Data access layer
│   │   ├── postgres/        # PostgreSQL implementation
│   │   ├── blockchain/      # Ethereum client
│   │   └── cache/           # Redis cache
│   ├── utils/               # Utility functions
│   └── contracts/           # Smart contract bindings
├── migrations/              # Database migrations
├── docs/                    # Swagger documentation
└── scripts/                 # Utility scripts
```

## API Endpoints

### AI Models

- `GET /api/v1/ai/models` - List all available models
- `POST /api/v1/ai/models` - Create a new model configuration
- `POST /api/v1/ai/chat` - Send chat completion request

### Documentation

- `GET /swagger/*` - Swagger UI

## Configuration

Key environment variables (see `.env.example` for full list):

```ini
# Server
STAGE_STATUS=dev              # dev or prod
SERVER_PORT=5000

# Database
DB_HOST=host.docker.internal
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=postgres

# Redis
REDIS_HOST=host.docker.internal
REDIS_PORT=6379

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_REFRESH_KEY=your-refresh-key

# Blockchain (optional)
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/YOUR-KEY
ETHEREUM_PRIVATE_KEY=your-private-key-hex
```

## Smart Contracts

Generate Go bindings from Solidity contracts:

```bash
# Compile contract
cd app/contracts
solc --abi --bin YourContract.sol -o build/

# Generate Go bindings
abigen --abi=build/YourContract.abi \
       --bin=build/YourContract.bin \
       --pkg=contracts \
       --out=your_contract.go
```

See `app/contracts/README.md` for detailed instructions.

## Testing & Quality

```bash
# Run full test suite with coverage
make test

# Code quality checks
make critic

# Security scan
make security

# Linting
make lint

# Build for production
make build
```

## Docker Commands

```bash
# Docker Compose
make compose.up         # Start all services
make compose.down       # Stop all services
make compose.logs       # View logs
make compose.restart    # Restart backend
make compose.rebuild    # Rebuild backend

# Individual services
make docker.postgres    # Start PostgreSQL
make docker.redis       # Start Redis
make docker.stop        # Stop all containers
```

## Architecture

Agent-C follows **Clean Architecture** principles:

- **Service Layer**: Business logic in `app/service/`
- **Repository Pattern**: Data access abstracted via `SqlStore` interface
- **Dependency Injection**: Dependencies injected through constructors
- **Middleware Chain**: JWT auth, logging, error handling

### Database Schema

- **providers** - AI model provider configurations
- **model_schemas** - JSON schemas for model options/responses
- **models** - Registry of available AI models
- **sellers** - API key providers (wallet-based)
- **consumers** - API key users (wallet-based)
- **api_keys** - Access keys with token tracking

## Adding Features

### New AI Provider

1. Create implementation in `app/service/model-providers/`
2. Define types in `app/types/`
3. Update routing in `app/service/consume.go`
4. Add provider to database

### New Endpoint

1. Add handler to `app/service/`
2. Register route in `app/routes/impl.go`
3. Add Swagger annotations
4. Run `make swag`

### Blockchain Integration

1. Write Solidity contract in `app/contracts/`
2. Compile and generate bindings
3. Use `EthereumClient` from `app/store/blockchain/`

## Error Handling

Standard JSON responses:

```json
// Error
{
  "error": true,
  "msg": "error message"
}

// Success
{
  "error": false,
  "msg": null,
  "data": {}
}
```

## Documentation

For detailed documentation, see:
- **[claude.md](./claude.md)** - Complete project context and architecture
- **[app/contracts/README.md](./app/contracts/README.md)** - Smart contract guide
- **Swagger UI** - http://localhost:8080/swagger/index.html (when running)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

Apache 2.0 © [William Bryce](https://github.com/wmbryce)

## Support

- Issues: [GitHub Issues](https://github.com/wmbryce/agent-c/issues)
- Documentation: See `claude.md` for comprehensive guide
