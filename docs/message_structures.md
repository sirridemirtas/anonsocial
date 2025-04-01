# Message and Notification API Structures

This document describes the data structures returned by the message and notification endpoints in the AnonSocial API.

## Message Endpoints

### GET `/api/v1/messages/unread-count`

Returns the count of unread messages for the authenticated user.

**Response Structure:**

```json
{
  "unreadCount": 5
}
```

### GET `/api/v1/messages`

Returns a list of conversations for the authenticated user.

**Response Structure:**

```json
[
  {
    "id": "60d21b4667d0d8992e610c85",
    "participants": [
      {
        "userId": "60d21b4667d0d8992e610c86",
        "username": "user1"
      },
      {
        "userId": "60d21b4667d0d8992e610c87",
        "username": "user2"
      }
    ],
    "updatedAt": "2023-06-01T10:00:00.000Z",
    "lastMessage": {
      "content": "Hello there!",
      "senderId": "60d21b4667d0d8992e610c87",
      "read": false,
      "timestamp": "2023-06-01T10:00:00.000Z"
    },
    "unreadCount": 1
  }
]
```

### GET `/api/v1/messages/{username}`

Returns a specific conversation with the specified user.

**Response Structure:**

```json
{
  "id": "60d21b4667d0d8992e610c85",
  "participants": [
    {
      "userId": "60d21b4667d0d8992e610c86",
      "username": "user1"
    },
    {
      "userId": "60d21b4667d0d8992e610c87",
      "username": "user2"
    }
  ],
  "messages": [
    {
      "id": "60d21b4667d0d8992e610c88",
      "senderId": "60d21b4667d0d8992e610c86",
      "content": "Hi there!",
      "read": true,
      "createdAt": "2023-06-01T09:00:00.000Z"
    },
    {
      "id": "60d21b4667d0d8992e610c89",
      "senderId": "60d21b4667d0d8992e610c87",
      "content": "Hello!",
      "read": false,
      "createdAt": "2023-06-01T10:00:00.000Z"
    }
  ]
}
```

### POST `/api/v1/messages/{username}`

Sends a message to a specified user.

**Request Body:**

```json
{
  "content": "Hello, how are you?"
}
```

**Response Structure:**
Returns the updated conversation object, identical to the GET `/api/v1/messages/{username}` response.

### POST `/api/v1/messages/{username}/read`

Marks all messages in a conversation as read.

**Response Structure:**

```json
{
  "message": "Mesajlar okundu olarak işaretlendi"
}
```

### DELETE `/api/v1/messages/{username}`

Deletes a conversation with the specified user.

**Response Structure:**

```json
{
  "message": "Görüşme silindi"
}
```

## Notification Endpoints

### GET `/api/v1/notifications/unread-count`

Returns the count of unread notifications for the authenticated user.

**Response Structure:**

```json
{
  "count": 5
}
```

### GET `/api/v1/notifications`

Returns a list of notifications for the authenticated user.

**Response Structure:**

```json
[
  {
    "id": "60d21b4667d0d8992e610c85",
    "type": "like",
    "senderUsername": "user2",
    "targetId": "60d21b4667d0d8992e610c86",
    "read": false,
    "createdAt": "2023-06-01T10:00:00.000Z"
  }
]
```

### PUT `/api/v1/notifications/{id}`

Marks a specific notification as read.

**Response Structure:**

```json
{
  "message": "Notification marked as read"
}
```

### PUT `/api/v1/notifications/mark-all-read`

Marks all notifications as read.

**Response Structure:**

```json
{
  "message": "All notifications marked as read"
}
```

### DELETE `/api/v1/notifications/delete-all`

Deletes all notifications.

**Response Structure:**

```json
{
  "message": "Tüm bildirimler silindi"
}
```

## Turkish Responses

Note that many API responses contain messages in Turkish. Below is a translation table for common messages:

| Turkish                              | English                     |
| ------------------------------------ | --------------------------- |
| "Mesajlar okundu olarak işaretlendi" | "Messages marked as read"   |
| "Görüşme silindi"                    | "Conversation deleted"      |
| "Tüm bildirimler silindi"            | "All notifications deleted" |
| "Çıkış başarılı"                     | "Logout successful"         |
| "Kullanıcı yetkisi güncellendi"      | "User role updated"         |
| "Avatar oluşturuldu"                 | "Avatar created"            |
| "Avatar güncellendi"                 | "Avatar updated"            |
