# Tarantool Key-Value Storage

![LICENSE-shield][license-shield-url] ![WORKFLOW-status][workflow-status-url]

A scalable **Key-Value storage** solution built on [Tarantool][1] and powered by [Gin (Go)][2]. This project offers CRUD operations, structured logging, and JSON-based communication.

## 🚀 Features

- **Tarantool Backend:** Leverages Tarantool for high performance and reliability.
- **RESTful API:** Provides standard CRUD endpoints for managing key-value pairs.
- **JSON-based Communication:** Accepts and returns JSON objects, enabling flexible data storage.
- **Structured Logging:** Built-in logging for easier debugging and traceability.
- **Graceful Error Handling:** Delivers clear HTTP status codes and detailed error messages.
- **Clean Architecture:** Modular and scalable design to support future growth.

## 🔧 Installation & Setup

### Prerequisites

- **Docker**: Ensure you have Docker installed on your system.
- **Docker Compose**: Required for orchestrating multiple containers.

### Quick Start

> [!IMPORTANT]
> If you are running application locally using `Docker Compose`, don't set `TT_HOST` environment variable to `localhost`.
> You can use any other arbitrary name. Compose will utilize built-in DNS resolver.
> [See compose documentation.][3]

1. **Clone the Repository:**

   ```sh
   git clone https://github.com/artem-burashniko/tarantool-app.git
   cd tarantool-app
   ```

2. **Configure Environment Variables**:

    Copy from `.env.example`:

    ```bash
    cp .env.example .env
    ```

    And set required fields. For instance:

    ```yaml
    # Application
    APP_NAME=tarantool-app
    APP_VERSION=0.1.0
    APP_ENV=local

    # Server
    HTTP_PORT=8080

    # Storage
    TT_HOST=tthost
    TT_PORT=3301
    TT_USER=user
    TT_PASSWORD=password
    ```

3. **Deploy with docker compose**:

    ```bash
    docker compose up -d
    ```

## 🌐 API Endpoints

Below are the available API endpoints for the Tarantool Key-Value Storage:

### 🔍 Get Value by Key

- **Description**: Retrieves the value for the specified key from the Tarantool database.
- **Method**: `GET`
- **Endpoint**: `/kv/{id}`
- **Path Parameters**:
  - `id` (string): The key ID to retrieve.
- **Responses**:
  - `200 OK`: Returns the key-value pair:

    ```json
    {
        "key": "foo",
        "value": {
            "bar": "baz"
        }
    }
    ```

  - `404 Not Found`: Key not found.
  - `500 Internal Server Error`: Server error.

---

### ➕ Create a New Key-Value Pair

- **Description**: Creates a new key-value pair in the Tarantool database.
- **Method**: `POST`
- **Endpoint**: `/kv`
- **Request Body**:

    ```json
    {
        "key": "foo",
        "value": {
            "bar": "baz"
        }
    }
    ```

- **Responses**:
  - `201 Created`: Key-value pair created successfully.

    ```json
    {
        "message": "created",
        "key": "foo",
        "value": {
            "bar": "baz"
        }
    }
    ```

  - `400 Bad Request`: Invalid request body or missing fields.
  - `409 Conflict`: Key already exists.
  - `500 Internal Server Error`: Server error.

---

### ✏️ Update Value by Key

- **Description**: Updates the value for the specified key in the Tarantool database.
- **Method**: `PUT`
- **Endpoint**: `/kv/{id}`
- **Path Parameters**:
  - `id` (string): The key ID to update.
- **Request Body**:

    ```json
    {
        "value": {
            "bar": "zab"
        }
    }
    ```

- **Responses**:
  - `200 OK`: Key-value pair updated successfully.

    ```json
    {
        "message": "updated",
        "key": "foo",
        "value": {
            "bar": "zab"
        }
    }
    ```

  - `400 Bad Request`: Invalid request body or missing fields.
  - `404 Not Found`: Key not found.
  - `500 Internal Server Error`: Server error.

---

### ❌ Delete Key-Value Pair

- **Description**: Deletes the specified key-value pair from the Tarantool database.
- **Method**: `DELETE`
- **Endpoint**: `/kv/{id}`
- **Path Parameters**:
  - `id` (string): The key ID to delete.
- **Responses**:
  - `200 OK`: Key-value pair deleted successfully.

    ```json
    {
        "message": "deleted",
        "key": "foo",
        "value": {
            "bar": "baz"
        }
    }
    ```

  - `404 Not Found`: Key not found.
  - `500 Internal Server Error`: Server error.

---

### 📘 Notes

- All endpoints accept and return JSON.
- Replace `{id}` with the actual key ID in the path.
- Ensure the Tarantool database is running and accessible before making requests.

## 📜 License

The project is licensed under an [MIT License](LICENSE).

<!-->
[license-shield-url]: https://img.shields.io/github/license/artem-burashnikov/tarantool-app?style=for-the-badge&color=blue
[workflow-status-url]: https://img.shields.io/github/actions/workflow/status/artem-burashnikov/tarantool-app/.github%2Fworkflows%2Fgo.yml?branch=main&style=for-the-badge

[1]: https://www.tarantool.io
[2]: https://github.com/gin-gonic/gin
[3]: https://docs.docker.com/reference/compose-file/services/#links
