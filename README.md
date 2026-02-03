# Bifrost

Bifrost is a Golang middleware API designed to route requests to specific databases and execute stored procedures or queries based on dynamic configuration. It acts as a secure gateway, managing authentication and database interactions efficiently.

## Features

- **Dynamic Database Routing**: Connects to differents databases dynamically based on the provided API Key.
- **Secure Authentication**: Validates requests using HMAC signatures (`X-SIGNATURE`), API Keys (`X-API-KEY`), and Authorization tokens.
- **Stored Procedure Execution**: Generic endpoint to execute stored procedures on the target database.
- **Safe Query Execution**: Generic endpoint to execute safe SQL queries.
- **Logging & Session Management**: Built-in request logging and session handling.
- **CORS Support**: Configured to handle Cross-Origin Resource Sharing.

## Prerequisites

- Go 1.22 or higher
- MySQL Database

## Installation

1.  Clone the repository:

    ```bash
    git clone <repository-url>
    cd api-mid-go
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

## Configuration

The application relies on environment variables for configuration. Create a `.env` file or set the following variable:

- `AES_ENC_KEY`: The key used for decrypting database passwords retrieved from the master database.

## Usage

Start the server:

```bash
go run main.go
```

The server will start on port `8081` (default).

### Endpoints

#### 1. Authorization Check

- **URL**: `/service/authorization`
- **Method**: `GET`
- **Description**: Checks if the request is authorized.

#### 2. Execute Stored Procedure

- **URL**: `/service/store-procedure`
- **Method**: `POST`
- **Description**: Executes a specific stored procedure.
- **Body**: JSON object containing `SpName` and `SpParams`.

#### 3. Execute Safe Query

- **URL**: `/service/safe-query`
- **Method**: `POST`
- **Description**: Executes a SQL query.

### Headers

All requests must include the following headers for authentication:

- `Content-Type`: `application/json`
- `X-API-KEY`: Your unique API Key.
- `Authorization`: Bearer token (or other auth token as configured).
- `X-SIGNATURE`: HMAC signature of the request body (if applicable).

## Project Structure

- `main.go`: Entry point of the application.
- `controllers/`: Contains the handler functions for the API endpoints.
- `utils/`: Utility functions for database connections, logging, encryption, etc.
- `models/`: Data structures and models.
