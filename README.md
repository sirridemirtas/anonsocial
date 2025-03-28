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
GO_ENV=production GIN_MODE=release go run main.go
```

## API Endpoints

Base URL: `/api/v1`

### Auth

- `POST /auth/register` - Register new user
- `POST /auth/login` - Login
  - Returns user information including username, email, and role
- `POST /auth/logout` - Logout (requires auth)
- `GET /auth/token-info` - Get information about the current token (requires auth)
- `POST /auth/refresh-token` - Refresh the authentication token (requires auth)
  - Generates a new token with a new expiration time
  - The old token is invalidated
  - Returns a success message: `{"message": "Token yenilendi"}`

### Users

- `GET /users` - Get all users
- `GET /users/:username` - Get user by ID
- `PUT /users/:id` - Update user (requires auth)
- `DELETE /users/:id` - Delete user (requires admin)
- `GET /auth/check-username/:username` - Check username
- `PUT /users/privacy` - (requires auth. isPrivate: true|false)
- `GET /users/:username/avatar` - Get user's avatar (respects privacy settings)
- `POST /users/:username/avatar` - Create or update user's avatar (requires auth, can only update own avatar)
- `POST /users/password/reset` - Reset password for the authenticated user (requires auth)
  - Example request:
    ```json
    {
      "currentPassword": "your-current-password",
      "newPassword": "your-new-password123"
    }
    ```

### Administration (Role 2 required)

- `PUT /admin/users/:username/role` - Update a user's role (can only set to 0 or 1)
  - Example request:
    ```json
    {
      "role": 1
    }
    ```
  - Response:
    ```json
    {
      "message": "Kullanıcı yetkisi güncellendi"
    }
    ```

### Feeds

- `GET /posts` - Get all posts (home feed)
- `GET /users/:username/posts` - Get posts by user (user feed)
- `GET /posts/:id/replies` - Get post replies (post feed)
- `GET /posts/university/:universityId` - Get posts by university (university feed)

### Posts

- `GET /posts/:id` - Get post
- `POST /posts` - Create new post (requires auth)
  - Request body:
    ```json
    {
      "content": "Your post content",
      "replyTo": "optional-post-id-to-reply-to",
      "universityId": "optional-university-id"
    }
    ```
  - If `universityId` is provided and valid, the post will be created for that university
  - If `universityId` is not provided, the post will be created for the user's own university
- `DELETE /posts/:id` - Delete post (requires auth)

### Replies

- `POST /posts/:id/replies` - Create reply (requires auth)
- `DELETE /posts/:id` - Delete post(reply) (requires auth)
  - Users with role 0 can only delete their own posts/replies
  - Users with role 1 or 2 can delete any post/reply regardless of ownership

### Reactions

- `POST /posts/:id/like` - Like post (requires auth)
- `POST /posts/:id/dislike` - Dislike post (requires auth)

- `POST /posts/:id/unlike` - Unlike post (requires auth)
- `POST /posts/:id/undislike` - Undislike post (requires auth)

### Messages

All message endpoints require authentication.

- `GET /messages` - Get list of all conversations for the authenticated user

  - Returns a summary of each conversation with just the last message

- `GET /messages/unread-count` - Get total number of unread messages across all conversations

  - Returns: `{"unreadCount": 5}`

- `GET /messages/:username` - Get conversation with specific user

  - Returns the conversation if it exists
  - Returns a 410 (Gone) if the conversation was deleted by the authenticated user
  - Returns a 400 (Bad Request) if the authenticated user tries to get a conversation with themselves

- `POST /messages/:username` - Send a message to a specific user

  - Request body: `{"content": "Message content"}`
  - Creates a new conversation if one doesn't exist
  - Returns a 400 (Bad Request) if the message content exceeds 500 characters
  - Returns a 400 (Bad Request) if the authenticated user tries to message themselves
  - Returns a 404 (Not Found) if the target user doesn't exist

- `POST /messages/:username/read` - Mark all messages in a conversation as read

  - Returns a 404 (Not Found) if the target user doesn't exist
  - Returns a 400 (Bad Request) if the conversation was deleted by the authenticated user

- `DELETE /messages/:username` - Delete conversation with specific user
  - Marks the conversation as deleted for the authenticated user only
  - The other participant can still see the conversation
  - Returns a 400 (Bad Request) if the conversation was already deleted by the authenticated user
  - Returns a 400 (Bad Request) if the authenticated user tries to delete a conversation with themselves

#### Message limits and behavior

- Each conversation stores a maximum of 100 messages
- When this limit is exceeded, the oldest messages are automatically removed
- Messages cannot be individually deleted, only entire conversations
- Deleted conversations are hidden from the user who deleted them but remain visible to the other user
- If both users delete a conversation, it is permanently removed from the database

### Notifications

All notification endpoints require authentication.

- `GET /notifications` - Get all notifications for the authenticated user

- `GET /notifications/unread-count` - Get the count of unread notifications

- `POST /notifications/:id` - Mark a notification as read

- `POST /notifications/mark-all-read` - Mark all notifications as read for the authenticated user

- `DELETE /notifications/delete-all` - Alternative endpoint to delete all notifications

### Contact

- `POST /contact` - Submit a contact form

  - Request body:
    ```json
    {
      "name": "Your Name",
      "email": "your.email@example.com",
      "subject": "Genel",
      "message": "Your message here"
    }
    ```
  - The `subject` field must be one of: "Genel", "Destek", "Öneri", "Teknik", "Şikayet"
  - Responses:
    - Success:
      ```json
      {
        "message": "Mesajınız gönderildi. En kısa sürede size dönüş yapacağız."
      }
      ```
    - Invalid email:
      ```json
      {
        "error": "Geçerli bir e-posta adresi girmelisiniz"
      }
      ```
    - Invalid subject:
      ```json
      {
        "error": "Konu, şu seçeneklerden biri olmalıdır: Genel, Destek, Öneri, Teknik, Şikayet"
      }
      ```

### Health

- `GET /health` - Check API health status
  - Returns basic system health information
  - No authentication required
  - Example response:
    ```json
    {
      "status": "up",
      "timestamp": "2023-05-01T12:34:56Z"
    }
    ```
