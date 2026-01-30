# Fraud Check Demo

A full-stack application demo that simulates a loan application fraud detection system while using FigChain to manage configuration updates in realtime.

## Overview

This demo consists of:
- **Backend**: Go service with REST API and SQLite database
- **Frontend**: React application for viewing loan applications

## Features

- Automatically seeds database with 50 fake loan applications
- Displays applicant names, credit scores, loan amounts, and risk factors
- Color-coded risk levels (Low, Medium, High)
- Application status tracking (Approved, Pending, Rejected)
- Filter applications by status
- Responsive design
- Long polling, for real-time updates

## Project Structure

```
fraudcheck/
├── backend/           # Go backend service
│   ├── main.go       # API server and routes
│   ├── models.go     # Data models
│   ├── database.go   # Database operations
│   ├── seed.go       # Fake data generation
│   └── go.mod        # Go dependencies
└── frontend/         # React frontend
    ├── src/
    │   ├── App.js    # Main application component
    │   ├── App.css   # Styles
    │   └── index.js  # Entry point
    ├── public/
    └── package.json  # NPM dependencies
```

## Prerequisites

- Go 1.25.5 or higher
- Node.js 16+ and npm
- Git

## Getting Started

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd demos/fraudcheck/backend
   ```

1. In the FigChain UI, navigate to Settings -> Client Credentials and create a new client credential. Choose the sandbox environment, and download the credential. Give it a `FailThreshold` field and make it an `int`. This serves as the threshold percentage in this example.

1. Create a schema in the namespace that you chose when creating a client credential, and call it `FraudCheckParameters`. Create a Fig in the same namespace, give it a key of `parameters` and add a version. Set whatever values you want.

1. Rename the downloaded client config file to `client-config.json` and place it in the `backend` directory that you cd'd into.

1. Install Go dependencies:
   ```bash
   go mod download
   ```

1. Run the backend server:
   ```bash
   go run .
   ```

The backend will:
- Start on `http://localhost:8080`
- Create a SQLite database (`fraudcheck.db`)
- Automatically seed 50 fake loan applications
- Fetch the initial configuration value from FigChain and listen for new values
- Expose API endpoints

### Frontend Setup

1. Open a new terminal and navigate to the frontend directory:
   ```bash
   cd demos/fraudcheck/frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The frontend will open in your browser at `http://localhost:3000`

## API Endpoints

### GET /api/applications
Retrieves all loan applications from the database.

**Response:**
```json
{
  "applications": [
    {
      "id": 1,
      "applicant_name": "John Smith",
      "credit_score": 720,
      "amount_requested": 50000,
      "risk_factor": 0.35,
      "status": "pending",
      "created_at": "2026-01-29T12:00:00Z"
    }
  ],
  "count": 50
}
```

### GET /api/health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

### GET /api/poll
Long polling endpoint for real-time updates.

**Response (when update available):**
```json
{
  "action": "refresh"
}
```

**Response (on timeout):**
```json
{
  "action": "none"
}
```

### POST /api/trigger-refresh
Manually trigger a refresh signal to all connected clients (useful for testing the long polling connection).

**Response:**
```json
{
  "message": "Refresh triggered"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/trigger-refresh
```

## Data Model

### Loan Application

| Field | Type | Description |
|-------|------|-------------|
| id | integer | Unique identifier |
| applicant_name | string | Full name of applicant |
| credit_score | integer | Credit score (300-850) |
| amount_requested | float | Loan amount requested |
| risk_factor | float | Calculated risk (0.0-1.0) |
| status | string | approved, pending, or rejected |
| created_at | datetime | Application timestamp |

## Risk Calculation

The risk factor is calculated based on:
- **Credit Score (70% weight)**: Lower scores increase risk
- **Loan Amount (30% weight)**: Higher amounts increase risk
- **Random variance**: ±10% for realistic variation

**Risk Levels:**
- Low: 0.0 - 0.3 (typically approved)
- Medium: 0.3 - 0.7 (pending review)
- High: 0.7 - 1.0 (typically rejected)

## Development

### Resetting the Database

To reset and reseed the database:

1. Stop the backend server
2. Delete the database file:
   ```bash
   rm backend/fraudcheck.db
   ```
3. Restart the backend server

### Modifying Seed Data

Edit [backend/seed.go](backend/seed.go) to:
- Change the number of generated applications
- Modify name lists
- Adjust risk calculation logic
- Customize loan amount ranges

## Technologies Used

**Backend:**
- Go 1.25.5
- Gin Web Framework
- SQLite3

**Frontend:**
- React 18
- CSS3
- Fetch API

## License

This is a demo application for educational purposes of FigChain.
