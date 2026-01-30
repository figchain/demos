# Fraud Check Frontend

React application for viewing and filtering loan applications.

## Quick Start

```bash
npm install
npm start
```

Opens at `http://localhost:3000`

## Features

- View all loan applications
- Filter by status (approved, pending, rejected)
- Color-coded risk levels
- Responsive table design

## Configuration

Backend API URL is configured in [App.js](src/App.js):
```javascript
const API_URL = 'http://localhost:8080/api';
```
