# AnonSocial Backend

A social media platform backend written in Go using Gin framework and MongoDB.

## Requirements

- Go 1.20+
- MongoDB

## Setup

1. Clone the repository:

```bash
git clone https://github.com/sirridemirtas/anonsocial.git
```

2. Install dependencies:

```bash
go mod download
```

3. Environment variables are already set in `.env.development`

## Running the Application

### Development Mode (with Hot Reload)

```bash
go run github.com/air-verse/air@latest
```

### Production Mode

```bash
GO_ENV=production go run main.go
```

## API Endpoints

### Auth

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout (requires auth)

### Users

- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user (requires auth)
- `DELETE /api/users/:id` - Delete user (requires admin)
