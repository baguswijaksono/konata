# Konata

Konata is a simple web application built with Go, Gin, and Gorm that allows you to execute curl commands and manage workspaces.

## Features

- Execute curl commands and save the history.
- Create and manage workspaces.

## Prerequisites

- Go 1.16 or later
- SQLite

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/konata.git
    cd konata
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

3. Build the application:

    ```sh
    go build -o konata main.go
    ```

4. Run the application:

    ```sh
    ./konata
    ```

    The server will start running at `http://localhost:8080`.

## API Endpoints

- `POST /execute`: Execute a curl command.
- `GET /history`: Retrieve the command execution history.
- `POST /workspace`: Create a new workspace.
- `GET /workspaces`: Retrieve all workspaces.

## Usage

You can use tools like `curl` or Postman to interact with the API endpoints.

### Execute a curl command

