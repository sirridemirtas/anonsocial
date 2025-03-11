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

Base URL: `/api/v1`

### Auth

- `POST /auth/register` - Register new user
- `POST /auth/login` - Login
- `POST /auth/logout` - Logout (requires auth)

### Users

- `GET /users` - Get all users
- `GET /users/:username` - Get user by ID
- `PUT /users/:id` - Update user (requires auth)
- `DELETE /users/:id` - Delete user (requires admin)
- `GET /auth/check-username/:username` - Check username

### Posts

- `GET /posts` - Get all posts
- `GET /posts/:id` - Get post by ID
- `GET /posts/university/:universityId` - Get posts by university
- `GET /posts/:id/replies` - Get post replies
- `POST /posts` - Create new post (requires auth)
- `DELETE /posts/:id` - Delete post (requires auth)
