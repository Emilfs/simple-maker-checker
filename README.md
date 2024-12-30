# Message Approval System

## Overview

The Message Approval System is a RESTful API developed in Go, designed to manage and process messages that require approval before being sent to their intended recipients. This system allows for the creation of message items, which can then be approved or rejected. Only approved messages are considered ready for sending.

## Design Decisions

- **Concurrent Safe Map**: Utilizes Go's `sync.Mutex` to ensure that operations on the shared `items` map are safe across multiple goroutines, preventing race conditions.
- **Status Constants**: Defined constants (`APPROVED`, `REJECTED`, `PENDING`) for clarity and to avoid hardcoding strings throughout the code.
- **Non-blocking Message Sending**: Utilizes Go's goroutines to simulate the message sending process in a non-blocking manner.
- **JSON Responses**: All API responses are in JSON format, providing a structured and easily parseable data format for clients.

## API

### Create Item

- **Method**: `POST`
- **URL**: `/items`
- **Body**:
  ```json
  {
    "id": "unique_item_id",
    "message": "Your message here",
    "recipient_id": "recipient_unique_id"
  }
  ```
- **Response**: The created item with a `pending` status.
- **Status Codes**: `201 Created` on success, `400 Bad Request` on invalid payload, `409 Conflict` if item already exists.

### Approve Item

- **Method**: `PUT`
- **URL**: `/items/{id}/approve`
- **Response**: The item with an `approved` status.
- **Status Codes**: `200 OK` on success, `404 Not Found` if item does not exist.

### Reject Item

- **Method**: `PUT`
- **URL**: `/items/{id}/reject`
- **Response**: The item with a `rejected` status.
- **Status Codes**: `200 OK` on success, `400 Bad Request` if item was already approved, `404 Not Found` if item does not exist.

### List Items

- **Method**: `GET`
- **URL**: `/items`
- **Response**: A list of all items.
- **Status Codes**: `200 OK` on success.

## How to Use

1. **Start the Server**: Run the Go application. The server will start listening on the specified port (default: 8080).

2. **Create an Item**: Send a `POST` request to `/items` with the item details in the request body.

3. **Approve/Reject an Item**: Send a `PUT` request to `/items/{id}/approve` or `/items/{id}/reject` to approve or reject an item, respectively.

4. **List Items**: Send a `GET` request to `/items` to retrieve a list of all items, including their statuses.

## Running the Application

Ensure you have Docker installed on your system. Navigate to the directory containing the application code and run:

```sh
docker build -t message-app .
docker run -d -p 8080:8080 message-app
```

The server will start, and you can interact with it using the endpoints described above through tools like `curl`, Postman, or any HTTP client of your choice.