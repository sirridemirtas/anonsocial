# AnonSocial API Testing Guide with Postman

This guide provides instructions for testing the AnonSocial API using the Postman collection.

## Prerequisites

- [Postman](https://www.postman.com/downloads/) installed on your machine
- Access to the AnonSocial API (local development or hosted environment)

## Setup

1. Import the `anonsocial_postman_collection.json` file into Postman
2. Create an environment in Postman with the following variables:
   - `baseUrl`: The base URL of the API (e.g., `http://localhost:8080/api/v1`)
   - `testUsername`: (Will be automatically generated during test run - alphanumeric only)
   - `testPassword`: (Will be automatically generated during test run)
   - `universityId`: A valid university ID (default is "115373")

> **Note:** Usernames in AnonSocial must contain only letters and numbers. The test scripts automatically generate compliant usernames.

## Running Tests

### Base URL Configuration

All requests in the collection use the `{{baseUrl}}` variable. Make sure this variable is properly set in your environment:

- Development: http://localhost:8080/api/v1

> **Note:** The API returns messages in Turkish. The test scripts have been updated to handle these responses.

### Sequence of Testing

For best results, run the endpoints in this order:

1. **Authentication (Phase 1)**

   - Register a new user (POST)
   - Login with the created user (POST)
   - Get token information (GET)

2. **User Profile Management**

   - Get user profile (GET)
   - Check username availability (GET)
   - Update user privacy (PUT)

3. **Avatar Management**

   - Update/Create user avatar (POST)
   - Get user avatar (GET)

4. **Post Management**

   - Create new post (POST)
   - Get post by ID (GET)
   - Create reply to post (POST)

5. **Post Interaction**

   - Delete reply (DELETE)
   - Delete post (DELETE)

6. **Password Management**

   - Reset user password (PUT)
   - Login with new password (POST)

7. **Authentication (Phase 2)**

   - Get token information (GET)

8. **Multi-User Interaction**

   - Register second user (POST)
   - Login as first user (POST)
   - Create post as first user (POST)
   - Logout first user (POST)
   - Login as second user (POST)
   - Get post as second user (GET)
   - Like post and reply to it (POST)
   - Send message to first user (POST)
   - Logout second user (POST)
   - Login first user again (POST)

9. **Notification Management**

   - Check unread notification count (GET)
   - Get all notifications (GET)
   - Mark single notification as read (PUT)
   - Mark all notifications as read (PUT)
   - Delete all notifications (DELETE)

10. **Message Management**

    - Get unread message count (GET)
    - Get conversation list (GET)
    - Get conversation with specific user (GET)
    - Mark conversation as read (POST)
    - Send reply message (POST)
    - Delete conversation (DELETE)

11. **Admin Operations** (requires admin privileges)

    - Update user role (PUT)

12. **Contact Form**

    - Submit contact form (POST)

13. **Final Cleanup**
    - Logout user (POST)
    - Delete test users (DELETE)

### Message API Response Format

When working with the messaging endpoints, note that the API returns:

- `/messages/unread-count` returns an object with an `unreadCount` property
- `/messages` returns an array of conversation objects with `id`, `participants`, and other properties
- `/messages/{username}` returns a single conversation object with `id`, `participants`, and `messages` properties
- When sending messages, the API returns the updated conversation object, not just a success message

### Automation

The collection includes test scripts to validate responses and automatically store important values like IDs and tokens for use in subsequent requests. These scripts run after each request to ensure proper testing flow.

Pre-request scripts generate random usernames and passwords for testing purposes, ensuring unique test data on each run.

## Expected Results

When run successfully, the collection should:

1. Create and authenticate users
2. Create, retrieve, and interact with posts
3. Test notification and messaging features
4. Verify proper error handling for invalid requests
5. Clean up test data to avoid cluttering the database

If any test fails, check the API logs for more detailed error information.
