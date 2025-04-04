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
make deps
```

3. Development mode (with hot reload):

```bash
make dev
```

4. Alternatively, build and run the application:

```bash
make run
```

This command first builds the application (which includes installing dependencies) and then runs it.

5. Or build and run separately:

```bash
make build  # Also installs dependencies
make clean  # Remove the compiled binary if needed
```

## Makefile Commands

- `make deps`: Downloads Go module dependencies
- `make dev`: Runs the application in development mode using Air for hot reloading
- `make build`: Builds the application with optimized flags for production
- `make clean`: Removes the compiled binary
- `make run`: Builds and runs the application in release mode

## Project Structure

```
anonsocial/
├── config/           # Application configuration management
├── controllers/      # HTTP request handlers and business logic
├── database/         # MongoDB connection and database operations
├── middleware/       # Gin middleware functions (auth, CORS, etc.)
├── models/           # Data models and structures
├── routes/           # API endpoint definitions and routing
├── utils/            # Helper functions and utilities
├── static/           # Static files (generated and served, gitignored)
├── .env.development  # Development environment variables
├── .env.production   # Production environment variables (gitignored)
├── main.go           # Application entry point
├── go.mod            # Go module definition
├── go.sum            # Go module checksumres
├── Makefile          # Build commands and development utilities
└── LICENSE           # MIT License

```

## Environment Variables

Create a `.env` file in the project root with the following variables:

```
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=anonsocial
JWT_SECRET=your_development_secret_key
JWT_EXPIRES_IN=720
COOKIE_DOMAIN=localhost
ALLOWED_ORIGINS=http://localhost:3000

```

- `PORT`: Server port (default: 8080)
- `MONGODB_URI`: MongoDB connection string
- `MONGODB_DB`: MongoDB database name
- `JWT_SECRET`: Secret key for JWT token generation
- `JWT_EXPIRES_IN`: JWT token expiration time in hours
- `COOKIE_DOMAIN`: Domain for authentication cookies
- `ALLOWED_ORIGINS`: CORS allowed origins (comma-separated)
- `GIN_MODE`: Gin framework mode (debug/release, set in Makefile)

# API Documentation

## Base URL

`/api/v1`

**Note:** Any path outside of `/api/v1` serves files from the `static` folder if they exist. This is typically used to serve frontend assets such as HTML, CSS, JavaScript files, and images stored in the `static/` directory.

## Authentication

Endpoints related to user authentication and token management.

| Method | Endpoint              | Parameters                                 | Description                                                                 |
| ------ | --------------------- | ------------------------------------------ | --------------------------------------------------------------------------- |
| POST   | `/auth/register`      | Body: `{username, password, universityId}` | Registers a new user.                                                       |
| POST   | `/auth/login`         | Body: `{username, password}`               | Authenticates the user and returns user information with token and cookie.  |
| POST   | `/auth/logout`        | None                                       | Logs out the current user (requires authentication) and deletes the cookie. |
| GET    | `/auth/token-info`    | None                                       | Retrieves information about the current token (requires auth).              |
| POST   | `/auth/refresh-token` | None                                       | Refreshes the authentication token (requires auth).                         |

- Most endpoints require authentication. The token obtained from `/auth/login` is sent via a cookie named `token`.

## User Management

Endpoints for managing user accounts and profiles.

| Method | Endpoint                           | Parameters                                        | Description                                                                                           |
| ------ | ---------------------------------- | ------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| GET    | `/users`                           | None                                              | Returns a list of all users (requires admin).                                                         |
| GET    | `/users/{username}`                | Path: username                                    | Retrieves details of a specific user.                                                                 |
| GET    | `/users/check-username/{username}` | Path: username                                    | Checks if a username is available.                                                                    |
| DELETE | `/users/{id}`                      | Path: id, Body: `{password}` (for self-deletion)  | Deletes a user account. Users can delete their own account with password; admins can delete any user. |
| PUT    | `/users/privacy`                   | Body: `{isPrivate: boolean}`                      | Updates the profile privacy setting (requires auth).                                                  |
| PUT    | `/users/password/reset`            | Body: `{currentPassword, newPassword}`            | Resets the password for the authenticated user (requires auth).                                       |
| GET    | `/users/{username}/avatar`         | Path: username                                    | Retrieves a user's avatar (respects privacy settings).                                                |
| POST   | `/users/{username}/avatar`         | Path: username, Body: JSON with avatar properties | Updates the user's avatar (requires auth, only own avatar).                                           |

## Posts

Endpoints for creating, retrieving, and interacting with posts.

| Method | Endpoint                | Parameters                                   | Description                                                                                                |
| ------ | ----------------------- | -------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| POST   | `/posts`                | Body: `{content, [universityId], [replyTo]}` | Creates a new post or reply. If replyTo is provided, it creates a reply to the specified post.             |
| GET    | `/posts/{id}`           | Path: id                                     | Retrieves a specific post.                                                                                 |
| GET    | `/posts/{id}/replies`   | Path: id, Query: `page=number`               | Retrieves replies to a specific post. Returns 50 replies per page.                                         |
| DELETE | `/posts/{id}`           | Path: id                                     | Deletes a post or reply. Users can delete their own content; moderators and admins can delete any content. |
| POST   | `/posts/{id}/like`      | Path: id                                     | Likes a post (requires auth).                                                                              |
| POST   | `/posts/{id}/dislike`   | Path: id                                     | Dislikes a post (requires auth).                                                                           |
| DELETE | `/posts/{id}/unlike`    | Path: id                                     | Removes a like from a post (requires auth).                                                                |
| DELETE | `/posts/{id}/undislike` | Path: id                                     | Removes a dislike from a post (requires auth).                                                             |

## Feeds

Endpoints for accessing different content feeds.

| Method | Endpoint                             | Parameters                               | Description                                                           |
| ------ | ------------------------------------ | ---------------------------------------- | --------------------------------------------------------------------- |
| GET    | `/feeds/home`                        | Query: `page=number`                     | Retrieves posts for the home feed. Returns 50 posts per page.         |
| GET    | `/feeds/universities/{universityId}` | Path: universityId, Query: `page=number` | Retrieves posts for a specific university. Returns 50 posts per page. |
| GET    | `/feeds/users/{username}`            | Path: username, Query: `page=number`     | Retrieves posts by a specific user. Returns 50 posts per page.        |

## Messages

Endpoints for private messaging between users.

| Method | Endpoint                    | Parameters                                | Description                                                                           |
| ------ | --------------------------- | ----------------------------------------- | ------------------------------------------------------------------------------------- |
| GET    | `/messages`                 | None                                      | Retrieves a list of all conversations for the authenticated user.                     |
| GET    | `/messages/{username}`      | Path: username                            | Retrieves the conversation with a specific user. Returns 410 if deleted, 400 if self. |
| POST   | `/messages/{username}`      | Path: username, Body: `{content: string}` | Sends a message to a specific user. Creates a new conversation if needed.             |
| DELETE | `/messages/{username}`      | Path: username                            | Deletes the conversation with a specific user (marks as deleted for the user).        |
| GET    | `/messages/unread-count`    | None                                      | Retrieves the total number of unread messages across all conversations.               |
| POST   | `/messages/{username}/read` | Path: username                            | Marks all messages in a conversation as read.                                         |

- Message conversations are limited to 100 messages; older messages are automatically deleted.

## Notifications

Endpoints for managing user notifications.

| Method | Endpoint                       | Parameters | Description                                                       |
| ------ | ------------------------------ | ---------- | ----------------------------------------------------------------- |
| GET    | `/notifications`               | None       | Retrieves all notifications for the user (last 50, unread first). |
| GET    | `/notifications/unread-count`  | None       | Retrieves the count of unread notifications.                      |
| PUT    | `/notifications/{id}`          | Path: id   | Marks a specific notification as read.                            |
| PUT    | `/notifications/mark-all-read` | None       | Marks all notifications as read.                                  |
| DELETE | `/notifications/delete-all`    | None       | Deletes all notifications.                                        |

- Notifications are limited to the last 50; older notifications are automatically deleted.

## Admin

Endpoints for administrative actions.

| Method | Endpoint                             | Parameters                          | Description                                           |
| ------ | ------------------------------------ | ----------------------------------- | ----------------------------------------------------- |
| PUT    | `/admin/users/{username}/role`       | Path: username, Body: `{role:0\|1}` | Updates a user's role (requires admin authorization). |
| PUT    | `/admin/users/{username}/activities` | None                                | Get user activities (requires admin authorization).   |

- User roles:
  - 0: Regular user
  - 1: Moderator
  - 2: Admin
- Certain actions require specific roles (e.g., deleting other users' posts, changing roles).

## Other Endpoints

Miscellaneous endpoints.

| Method | Endpoint   | Parameters                              | Description                                                                                      |
| ------ | ---------- | --------------------------------------- | ------------------------------------------------------------------------------------------------ |
| POST   | `/contact` | Body: `{name, email, subject, message}` | Submits a contact form. Subject must be one of: "Genel", "Destek", "Öneri", "Teknik", "Şikayet". |
| GET    | `/health`  | None                                    | Checks the API's health status.                                                                  |
