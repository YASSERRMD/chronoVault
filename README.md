# ChronoVault

Time-based contract obligation tracking system.

## Quick Start

```bash
# Using Docker Compose
docker-compose up -d

# Or run locally
cd backend && go run cmd/server/main.go
cd frontend && npm install && npm run dev
```

## Default Credentials

- Email: admin@demo.com  
- Password: password123

## Tech Stack

- **Frontend**: Vue 3 + Vite
- **Backend**: Go + Gin
- **Database**: SQLite
- **Real-time**: WebSocket
- **Containerized**: Docker

## API Endpoints

- `POST /api/auth/login` - Login
- `POST /api/auth/register` - Register
- `GET /api/contracts` - List contracts
- `GET /api/obligations` - List obligations
- `GET /api/reports/financial-summary` - Financial summary
- `GET /api/reports/risk-exposure` - Risk exposure

## WebSocket

Connect to `ws://localhost:8080/ws` for real-time updates.

Events:
- `obligation_activated`
- `obligation_breached`
- `obligation_fulfilled`
- `penalty_applied`
