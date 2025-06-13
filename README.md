# Balance Transactions API

A REST API service for processing user balance transactions built with Go, Echo framework, and PostgreSQL.

## Setup and Running

### Prerequisites
- Docker and Docker Compose

### Quick Start
```bash
docker compose up -d
```

This will start:
- PostgreSQL database on port 5432
- API server on port 8080

The database will be initialized with 3 test users (IDs: 1, 2, 3) with zero balance.

## API Endpoints

### POST /user/{userId}/transaction
Process a transaction for a user.

**Headers:**
- `Source-Type: game|server|payment`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "state": "win",
  "amount": "10.15",
  "transactionId": "unique-transaction-id"
}
```

**Response:**
- `200 OK` - Transaction processed successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - User not found
- `409 Conflict` - Duplicate transaction ID

### GET /user/{userId}/balance
Get current user balance.

**Response:**
```json
{
  "userId": 1,
  "balance": "10.15"
}
```

## Testing

### Add Balance (Win Transaction)
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{"state": "win", "amount": "50.00", "transactionId": "tx-001"}'
```

### Subtract Balance (Lose Transaction)
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{"state": "lose", "amount": "10.50", "transactionId": "tx-002"}'
```

### Get Balance
```bash
curl http://localhost:8080/user/1/balance
```

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/service -v
go test ./internal/handler -v
```

### Stop Services
```bash
docker compose down
```

### View Logs
```bash
docker compose logs -f
```

### Reset Database
```bash
docker compose down -v
docker compose up -d
```
