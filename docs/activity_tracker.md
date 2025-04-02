# User Activity Tracker Middleware

## Overview

The Activity Tracker is a Gin middleware that records user activity on specific endpoints. It tracks the following information:

- Username (extracted from JWT token)
- Client's real IP address (using X-Forwarded-For header or RemoteAddr)
- Port number (parsed from RemoteAddr)
- Endpoint name (e.g., '/posts' or '/messages/:username')
- HTTP method (e.g., 'POST', 'GET')
- Request timestamp in ISO 8601 format (UTC)

This middleware is designed to be lightweight, idempotent, and fault-tolerant, making it suitable for REST APIs.

## Architecture

The middleware follows these principles:

1. **One document per user**: Each user has a single document in the MongoDB collection.
2. **IP-Port tracking**: Each document contains an array of IP-Port combinations used by the user.
3. **Action history**: Each IP-Port entry contains an array of actions performed from that IP-Port.
4. **Efficiency**: If an IP-Port combination already exists, only a new action is added to the existing entry.

## MongoDB Document Structure

Below is an example of how the documents look in the MongoDB collection:

```json
{
  "username": "user123",
  "ipEntries": [
    {
      "ip": "203.0.113.1",
      "port": "54321",
      "actions": [
        {
          "endpoint": "/api/v1/posts",
          "method": "POST",
          "timestamp": "2025-03-10T15:30:45.123Z"
        },
        {
          "endpoint": "/api/v1/posts",
          "method": "POST",
          "timestamp": "2025-03-10T16:12:33.456Z"
        }
      ],
      "updatedAt": "2025-03-10T16:12:33.456Z"
    },
    {
      "ip": "198.51.100.42",
      "port": "39876",
      "actions": [
        {
          "endpoint": "/api/v1/messages/otheruser",
          "method": "POST",
          "timestamp": "2025-03-11T09:45:22.789Z"
        }
      ],
      "updatedAt": "2025-03-11T09:45:22.789Z"
    }
  ],
  "updatedAt": "2025-03-11T09:45:22.789Z"
}
```

## IP Address Detection

The middleware tries to get the most accurate client IP by:

1. First checking the X-Forwarded-For header (for clients behind proxies)
2. Then checking the X-Real-IP header (another common proxy header)
3. Finally falling back to the RemoteAddr from the HTTP request

This ensures accurate IP detection even when the client is behind a proxy or NAT.

## Usage

To apply the middleware to an endpoint:

```go
router.POST("/some-endpoint", middleware.Auth(0), middleware.ActivityTracker(), controllers.SomeController)
```

The middleware should be applied after authentication middleware to ensure the username is available.

## Admin API Endpoints

### Get User Activities

```
GET /api/v1/admin/users/:username/activities
```

This endpoint allows administrators to retrieve all activity records for a specific user, including all IP addresses, ports, and actions performed.

**Required Permission**: Administrator Role (2)

**Parameters**:

- `username` - The username of the user whose activities to retrieve

**Response**: Returns the complete user activity document with all IP entries and associated actions.

## Implementation Details

- MongoDB indexes are created on the username field for faster lookups
- The middleware executes after the request is processed (c.Next() is called first)
- If the username cannot be determined, the activity is not recorded
- Error handling ensures that request processing continues even if activity tracking fails
