# Agent-C Project Context

## Project Overview

Agent-C is a Go-based AI model aggregation and blockchain integration platform. It provides a unified API for interacting with multiple AI model providers while managing access through blockchain-based authentication and payment systems.

## Tech Stack

- **Language**: Go 1.25.4
- **Web Framework**: Fiber v2 (Express.js-inspired Go framework)
- **Database**: PostgreSQL with pgx/v5 driver
- **Cache**: Redis
- **Blockchain**: Ethereum integration via go-ethereum
- **AI Integration**: OpenAI API via sashabaranov/go-openai
- **Authentication**: JWT tokens via golang-jwt/jwt
- **API Documentation**: Scalar/OpenAPI via swaggo/swag
- **Validation**: go-playground/validator
- **Logging**: zerolog
- **Migrations**: Goose

## Architecture Pattern

The project follows **Clean Architecture** principles with clear separation of concerns:

### Directory Structure

```
agent-c/
├── cmd/                      # Application entry point and configs
│   ├── main.go              # Main application entry
│   └── configs/             # Server configurations
│       └── fiber_config.go
├── app/                      # Business logic layer
│   ├── middleware/          # HTTP middleware
│   │   ├── fiber_middleware.go
│   │   └── jwt_middleware.go
│   ├── types/               # Domain models and DTOs
│   │   ├── chat_model.go    # AI chat-related types
│   │   ├── openai_model.go  # OpenAI-specific types
│   │   └── contract_model.go # Blockchain contract types
│   ├── service/             # Business logic services
│   │   ├── impl.go          # Service implementation
│   │   ├── models.go        # Model management
│   │   ├── consume.go       # AI model consumption
│   │   └── model-providers/ # AI provider integrations
│   │       └── openai.go
│   ├── routes/              # HTTP route definitions
│   │   └── impl.go
│   ├── store/               # Data access layer abstraction
│   │   ├── impl.go          # Store interface
│   │   ├── postgres/        # PostgreSQL implementation
│   │   │   ├── impl.go
│   │   │   └── models.go
│   │   ├── blockchain/      # Ethereum client
│   │   │   └── ethereum.go
│   │   └── cache/           # Redis cache
│   ├── utils/               # Utility functions
│   │   ├── jwt_generator.go
│   │   ├── jwt_parser.go
│   │   ├── password_generator.go
│   │   ├── validator.go
│   │   ├── connection_url_builder.go
│   │   ├── start_server.go
│   │   └── contract_utils.go
│   └── contracts/           # Smart contract bindings
│       ├── README.md
│       └── ERC20.sol
├── migrations/              # Database migrations
│   └── 20251204040426_init_tables.sql
├── docs/                    # Scalar API documentation
│   └── swagger.json
└── scripts/                 # Utility scripts

```

## Core Features

### 1. AI Model Management

The platform provides a unified interface for managing and consuming AI models from various providers:

- **Model Registry**: Store model configurations with provider details, schemas, and endpoints
- **Provider Abstraction**: Support for multiple AI providers (currently OpenAI)
- **Flexible Schema**: JSON-based schema definitions for model options and responses
- **Unified API**: Single endpoint (`/api/v1/ai/chat`) for all model interactions

### 2. Blockchain Integration

Ethereum blockchain integration for decentralized access control and payments:

- **Smart Contract Support**: Generate Go bindings from Solidity contracts
- **ERC20 Token Integration**: Built-in support for token operations
- **Wallet Authentication**: Support for wallet-based authentication
- **Transaction Management**: Handle contract calls and state changes
- **Multi-role System**: Sellers (API key providers) and Consumers (API users)

### 3. Access Control & Payments

Token-based access system:

- **API Key Management**: Sellers can create API keys with token allocations
- **Token Accounting**: Track token usage per API key
- **Provider Mapping**: Link API keys to specific AI providers
- **JWT Authentication**: Secure endpoint access

## Database Schema

The application uses a PostgreSQL database with schema `agc`:

### Tables

1. **providers**
   - Stores AI model provider information
   - Fields: id (UUID), name, description, endpoint_url, timestamps

2. **model_schemas**
   - JSON schemas for model options and responses
   - Fields: id (UUID), type (options/response), name, schema (JSONB), timestamps
   - Types: 'options' for input schemas, 'response' for output schemas

3. **models**
   - Registry of available AI models
   - Fields: id (UUID), model_key (unique), name, description, provider_id, options_schema_id, response_schema_id, request_url, timestamps
   - References: providers, model_schemas (2x)

4. **sellers**
   - Entities that provide API keys
   - Fields: id (UUID), wallet_address (unique), timestamps

5. **consumers**
   - Entities that use API keys
   - Fields: id (UUID), wallet_address (unique), timestamps

6. **api_keys**
   - API key management with token tracking
   - Fields: id (UUID), api_key (unique), tokens_available, provider_id, seller_id, timestamps
   - References: providers, sellers

## API Endpoints

### AI Model Endpoints

- `GET /api/v1/ai/models` - List all available models
- `POST /api/v1/ai/models` - Create a new model configuration
- `POST /api/v1/ai/chat` - Send chat completion request to AI model

### Documentation

- `GET /docs/* - Scalar API Reference UI

## Environment Configuration

The application requires the following environment variables:

### Server Settings
- `STAGE_STATUS` - "dev" or "prod" (controls graceful shutdown)
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_READ_TIMEOUT` - Read timeout in seconds

### Database Settings (PostgreSQL)
- `DB_TYPE` - Database type (pgx or mysql)
- `DB_HOST` - Database host
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSL_MODE` - SSL mode (disable, require, etc.)
- `DB_MAX_CONNECTIONS` - Max connection pool size
- `DB_MAX_IDLE_CONNECTIONS` - Max idle connections
- `DB_MAX_LIFETIME_CONNECTIONS` - Connection lifetime

### Redis Settings
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password (optional)
- `REDIS_DB_NUMBER` - Redis database number

### JWT Settings
- `JWT_SECRET_KEY` - Secret key for access tokens
- `JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT` - Access token expiry (minutes)
- `JWT_REFRESH_KEY` - Secret key for refresh tokens
- `JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT` - Refresh token expiry (hours)

### Blockchain Settings (Optional)
- `ETHEREUM_RPC_URL` - Ethereum node RPC URL
- `ETHEREUM_PRIVATE_KEY` - Private key for signing transactions (hex format, without 0x)

### Migration Settings (Goose)
- `GOOSE_DRIVER` - Migration driver (postgres)
- `GOOSE_DBSTRING` - Database connection string for migrations
- `GOOSE_MIGRATION_DIR` - Migration files directory

## Development Setup

### Prerequisites
- Go 1.25.4+
- Docker & Docker Compose
- PostgreSQL 14+
- Redis 7+
- Solidity compiler (for smart contracts)
- abigen (for contract bindings)

### Quick Start

1. **Using Docker Compose (Recommended)**
   ```bash
   # Start all services
   make compose.up
   
   # View logs
   make compose.logs
   
   # Restart backend
   make compose.restart
   
   # Stop all services
   make compose.down
   ```

2. **Local Development with Air (Hot Reload)**
   ```bash
   # Install Air
   make air.install
   
   # Start with hot reload
   make dev
   ```

3. **Manual Setup**
   ```bash
   # Start infrastructure
   make docker.postgres
   make docker.redis
   
   # Run migrations
   make goose.up
   
   # Generate API docs
   make swag
   
   # Run application
   go run cmd/main.go
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

# Reset database (careful!)
make goose.reset
```

### Smart Contract Development

1. Write Solidity contract in `app/contracts/`
2. Compile contract:
   ```bash
   solc --abi --bin YourContract.sol -o build/
   ```
3. Generate Go bindings:
   ```bash
   abigen --abi=build/YourContract.abi \
          --bin=build/YourContract.bin \
          --pkg=contracts \
          --out=../app/contracts/your_contract.go
   ```

## Key Design Patterns

### Service Layer Pattern
- `Service` struct acts as the main business logic coordinator
- Injected with logger, store, and fiber app dependencies
- Methods attached to Service handle HTTP requests and responses

### Repository Pattern
- `SqlStore` interface abstracts data access
- Multiple implementations: PostgreSQL, (extensible to other databases)
- Clean separation between business logic and data access

### Dependency Injection
- Dependencies injected through constructors
- Main function orchestrates initialization:
  1. Logger setup
  2. Database connection
  3. Fiber app initialization
  4. Service creation with dependencies
  5. Route registration

### Middleware Chain
- JWT authentication middleware
- Request logging
- Error handling
- CORS configuration

## Testing & Quality

```bash
# Run full test suite with coverage
make test

# Individual checks
make critic    # Code quality checks
make security  # Security vulnerability scan
make lint      # Linter checks
```

## Deployment

### Docker Build
```bash
# Production build
docker build -t agent-c .

# Development build (with hot reload)
docker build -f Dockerfile.dev -t agent-c-dev .
```

### Docker Compose Services
- `postgres` - PostgreSQL database (port 5432)
- `redis` - Redis cache (port 6379)
- `backend` - Go application (port 8080)

All services connected via `dev-network` bridge network with health checks.

## API Documentation

Scalar API Reference available at: `http://localhost:8080/docs`

Auto-generated from code annotations using swaggo/swag.

## Common Tasks

### Adding a New AI Provider

1. Create provider implementation in `app/service/model-providers/`
2. Define provider-specific types in `app/types/`
3. Update `ConsumeModel` function to route to provider
4. Add provider configuration to database

### Adding New Endpoints

1. Define handler in `app/service/`
2. Add route in `app/routes/impl.go`
3. Add Swagger annotations
4. Run `make swag` to update documentation

### Adding Blockchain Features

1. Write Solidity contract in `app/contracts/`
2. Compile and generate Go bindings
3. Create service methods in `app/service/contract.go`
4. Use `EthereumClient` from `app/store/blockchain/`

## Error Handling

Standard error response format:
```json
{
  "error": true,
  "msg": "error message or validation errors"
}
```

Success response format:
```json
{
  "error": false,
  "msg": "success message or null",
  "data": {}
}
```

## Logging

Uses structured logging with zerolog:
- Console output in development
- Timestamp format: Unix timestamp
- Context-aware logging throughout the application

## Security Considerations

- JWT-based authentication
- Password hashing with bcrypt
- Ethereum private key management
- SQL injection prevention via parameterized queries (pgx)
- Request validation via go-playground/validator
- CORS configuration
- API rate limiting (TODO)

## Future Enhancements

- Additional AI provider integrations (Anthropic, Cohere, etc.)
- Smart contract-based payment settlements
- Token staking mechanisms
- Multi-chain support
- Enhanced monitoring and metrics
- API rate limiting and quotas
- WebSocket support for streaming responses
- Circuit breaker pattern for external services

## Project Status

Based on recent commits:
- Clean architecture refactoring completed
- Docker Compose configuration updated
- Core infrastructure in place
- Active development phase
