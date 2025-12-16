# Soccer Manager API

A RESTful API for managing fantasy football teams, built with Go using hexagonal architecture.

## Features

- User authentication (JWT-based)
- Team management (create, view, update)
- Player management (view, update player information)
- Transfer market (list players, buy/sell players)
- Redis caching for improved performance
- Localization support (English and Georgian)
- PostgreSQL database
- Docker support

## Architecture

The project follows **Hexagonal Architecture (Ports & Adapters)** pattern:

- **Domain Layer**: Core business entities and rules
- **Ports Layer**: Interfaces for external dependencies
- **Application Layer**: Business logic and use cases
- **Infrastructure Layer**: Database, cache, HTTP handlers

## Prerequisites

- Docker and Docker Compose

## Setup

1. Clone the repository:
```bash
git clone https://github.com/temo927/Soccer-Manager-API.git
cd Soccer-Manager-API
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Start all services (PostgreSQL, Redis, and API):
```bash
cd docker
docker compose up -d
```

4. The API will be available at `http://localhost:8080`

### Stopping Services

To stop all services:
```bash
cd docker
docker compose down
```

To stop and remove volumes (clears database):
```bash
cd docker
docker compose down -v
```

### Viewing Logs

View API logs:
```bash
cd docker
docker compose logs api
```

View all logs:
```bash
cd docker
docker compose logs
```

Follow logs in real-time:
```bash
cd docker
docker compose logs -f api
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user

### Team Management
- `GET /api/v1/teams/me` - Get current user's team
- `PUT /api/v1/teams/me` - Update team name/country
- `GET /api/v1/teams/me/players` - Get team's players

### Player Management
- `GET /api/v1/players/{id}` - Get player details
- `PUT /api/v1/players/{id}` - Update player (first_name, last_name, country)

### Transfer List
- `POST /api/v1/players/{id}/transfer-list` - List player for transfer
- `DELETE /api/v1/players/{id}/transfer-list` - Remove player from transfer list
- `GET /api/v1/transfer-list` - Get all players on transfer list
- `POST /api/v1/transfer-list/{listing_id}/buy` - Buy player from transfer list

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Localization

The API supports English (en) and Georgian (ka) languages. Set the `Accept-Language` header:
```
Accept-Language: en
```
or
```
Accept-Language: ka
```

## Response Format

All responses follow this format:
```json
{
  "success": true,
  "data": {},
  "message": "Success message",
  "errors": []
}
```

## Testing

### Unit Tests
Run unit tests:
```bash
go test ./...
```

### Integration Tests
Integration tests are located in `tests/integration/`. Run them with:
```bash
go test ./tests/integration/...
```

## Project Structure

```
Soccer-Manager-API/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/          # Domain entities and business rules
│   ├── ports/          # Interfaces (repositories, cache)
│   ├── app/            # Use cases and business logic
│   └── infrastructure/ # External adapters (DB, HTTP, cache)
├── pkg/                # Shared utilities (JWT, password, localization)
├── tests/              # Test files
├── api/postman/        # Postman collection
└── docker/             # Docker configuration
```

## Environment Variables

See `.env.example` for all available environment variables.


