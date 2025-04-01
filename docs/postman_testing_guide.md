# AnonSocial API Testing Guide with Postman

This guide explains how to use the Postman collection for testing the AnonSocial API endpoints.

## Setup

### Prerequisites

1. [Postman](https://www.postman.com/downloads/) installed
2. AnonSocial server running (default: http://localhost:8080)

### Import the Collection

1. Open Postman
2. Click "Import" button in the top left
3. Select the `anonsocial_postman_collection.json` file
4. The collection will appear in your Postman workspace

### Set Up Environment

1. Create a new Environment in Postman by clicking on "Environments" in the sidebar and then "+" button
2. Name it "AnonSocial Development"
3. Add the following variables:
   - `baseUrl`: http://localhost:8080/api/v1 (or your server URL)
   - `testUsername`: Leave empty (will be auto-generated)
   - `testPassword`: Leave empty (will be auto-generated)
   - `universityId`: 115373 (this is a valid university ID, do not change)
   - `token`: Leave empty (will be set during tests)
   - `testUserId`: (Optional) ID of an existing user for testing
   - `postId`: (Optional) ID of an existing post for testing
   - `targetUsername`: (Optional) Username of another user to test admin operations
4. Save the environment
5. Select the environment from the dropdown in the top right corner

## Running Tests

### Base URL Configuration

All requests in the collection use the `{{baseUrl}}` variable. Make sure this variable is properly set in your environment:

- Development: http://localhost:8080
- Staging: https://staging-api.anonsocial.com
- Production: https://api.anonsocial.com

### Sequence of Testing

For best results, run the endpoints in this order:

1. **Authentication (Phase 1)**

   - Register a new user
   - Login with the created user
   - Get token information

2. **User Profile Management**

   - Get user profile
   - Check username availability
   - Update user privacy

3. **Avatar Management**

   - Update/Create user avatar (first)
   - Get user avatar (after creating it)

4. **Post Management**

   - Create new post
   - Get post by ID
   - Create reply to post

5. **Post Interaction**

   - Delete reply (uses the regular post deletion endpoint)
   - Delete post

6. **Password Management**

   - Reset user password
   - Login with new password (required after password reset)

7. **Authentication (Phase 2)**

   - Get token information (verify new login)

8. **Admin Operations** (requires admin privileges)

   - Update user role

9. **Cleanup**
   - Logout (do this last)
   - Delete user (requires moderator role, very last step)

### Automation

You can run the entire collection as an automated test:

1. Click the three dots (...) next to the collection name
2. Select "Run collection"
3. Configure the run settings
4. Click "Run AnonSocial API"

## Environment Variables

The collection uses pre-request scripts to automatically generate test values:

- Random usernames
- Random passwords
- Valid university ID (preset to "115373")

It also automatically captures authentication tokens from responses and stores them for subsequent requests.

### Important Note About University ID

The `universityId` is NOT randomly generated because it must be a valid ID from the system's recognized university list. The collection uses "115373" as a valid ID for testing purposes. If this ID becomes invalid, you'll need to manually update it in your environment with another valid university ID.

## Troubleshooting

### Authentication Issues

- If you see 401 Unauthorized errors, make sure your token is valid
- The collection automatically captures tokens from login responses
- You may need to re-login if your session expires

### Missing Resources

- Some tests require existing resources (post IDs, user IDs)
- If you don't have these resources, add them manually to your environment

### Base URL Problems

- If you're getting connection errors, verify that the `baseUrl` variable is set correctly for your environment

## Notes for Developers

- The test scripts include proper validation for different response codes
- Failed tests include descriptive messages to help identify issues
- The collection handles both success and error paths for most endpoints
