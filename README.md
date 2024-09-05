# User Management Service

A simple microservice implemented in Go that provides an HTTP API for managing user data. The service uses an in-memory storage mechanism to store user data.

## Features

- Add a new user
- Update an existing user
- Delete a user
- Return a paginated list of users with filtering options
- Health check endpoint
- In-memory data storage

## Requirements

- Go 1.20 or higher
- Docker (if you want to run the application in a Docker container)

## Getting Started

### Running the Application Locally

To run the application locally, follow these steps:

1. **Clone the repository**:

   ```sh
   git clone <repository-url>
   cd user-management-service

2. **Build the application**:
   ```sh
   go build -o user-management-service .

3. **Run the application**:
   ```sh
   ./user-management-service

### Running the Application with Docker

1. **Build the Docker image**:
   ```sh
   docker build -t user-management-service .

2. **Run the Docker container**:
   ```sh
   docker run --name my-user-service -p 8080:8080 user-management-service


### API Endpoints

1. **Add a New User**:
**Endpoint**: /users

Method: POST

Request Body:
```json
{
    "first_name": "Alice",
    "last_name": "Smith",
    "nickname": "Ali123",
    "password": "securepassword",
    "email": "alice@example.com",
    "country": "UK"
}
```
Response:

Status Code: 201 Created
Body:
```json
{
    "id": "generated-uuid",
    "first_name": "Alice",
    "last_name": "Smith",
    "nickname": "Ali123",
    "email": "alice@example.com",
    "country": "UK",
    "created_at": "2023-09-01T12:34:56Z",
    "updated_at": "2023-09-01T12:34:56Z"
}
```

2. **Update an Existing User** :

Endpoint: /users/{id}

Method: PUT

Request Body:
```json

{
    "first_name": "Alice",
    "last_name": "Smithers",
    "nickname": "Ali456",
    "email": "alice.smith@example.com",
    "country": "US"
}

```
Response:

Status Code: 200 OK

Body
``` json
{
    "id": "user-id",
    "first_name": "Alice",
    "last_name": "Smithers",
    "nickname": "Ali456",
    "email": "alice.smith@example.com",
    "country": "US",
    "created_at": "2023-09-01T12:34:56Z",
    "updated_at": "2023-09-01T13:00:00Z"
}
```

3. **Delete a User**
Endpoint: /users/{id}

Method: DELETE

Response:

Status Code: 200

```json
{
    "message": "User deleted successfully"
}
```

4. **Get a Paginated List of Users**
Endpoint: /users

Method: GET

Query Parameters:

page (optional): Page number (default is 1)
pageSize (optional): Number of users per page (default is 10)
country (optional): Filter users by country
Example Request:
```sh
GET /users?page=1&pageSize=2&country=UK
```
Response:

Status Code: 200 OK
Body:
``` json
{
    "total": 3,
    "users": [
        {
            "id": "user-id-1",
            "first_name": "Alice",
            "last_name": "Smith",
            "nickname": "Ali123",
            "email": "alice@example.com",
            "country": "UK",
            "created_at": "2023-09-01T12:34:56Z",
            "updated_at": "2023-09-01T12:34:56Z"
        },
        {
            "id": "user-id-2",
            "first_name": "Bob",
            "last_name": "Johnson",
            "nickname": "BobbyJ",
            "email": "bob@example.com",
            "country": "UK",
            "created_at": "2023-09-01T12:35:56Z",
            "updated_at": "2023-09-01T12:35:56Z"
        }
    ]
}
```

5. **Health Check**:
Endpoint: /health

Method: GET

Response:
Status Code: 200 OK
Body:

```json
{
    "status": "OK"
}
```