# Language Learning Portal Backend

This is the backend server for a language learning portal that serves as:
- An inventory of vocabulary
- A Learning Record Store (LRS)
- A unified launchpad for learning apps

## Tech Stack

- Go
- SQLite3
- Gin (Web Framework)
- Mage (Task Runner)

## Project Structure

```
backend_go/
├── cmd/
│   └── server/      # Main application entry point
├── internal/
│   ├── models/      # Data structures and database operations
│   ├── handlers/    # HTTP handlers
│   └── service/     # Business logic
├── db/
│   ├── migrations/  # Database migrations
│   └── seeds/       # Seed data
├── magefile.go      # Task runner definitions
├── go.mod          # Go module file
└── words.db        # SQLite database
```

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Install Mage:
```bash
go install github.com/magefile/mage
```

3. Initialize the database:
```bash
mage initdb
```

4. Run migrations:
```bash
mage migrate
```

5. Seed the database:
```bash
mage seed
```

## Development

To run the server:
```bash
go run cmd/server/main.go
```

## API Documentation

The API provides endpoints for:
- Dashboard statistics
- Word management
- Group management
- Study session tracking
- System reset

For detailed API documentation, see the API specification in the technical documentation.
