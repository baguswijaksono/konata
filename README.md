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
    git clone https://github.com/baguswijaksono/re.git
    cd re
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

3. Build the application:

    ```sh
    go build -o re main.go
    ```

4. Run the application:

    ```sh
    ./re
    ```

    The server will start running at `http://localhost:8080`.
