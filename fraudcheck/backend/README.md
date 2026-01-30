# Fraud Check Backend

Go REST API service for managing loan application data.

## Quick Start

```bash
go mod download
go run .
```

Server starts on `http://localhost:8080`

## Endpoints

- `GET /api/applications` - Get all loan applications
- `GET /api/health` - Health check

## Database

Uses SQLite3 (`fraudcheck.db`) with automatic schema creation and seeding.
