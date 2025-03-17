# Order Service

This project is an order service application written in Go. It provides APIs to create and retrieve orders. The application uses PostgreSQL as its database and is containerized using Docker.

## Project Structure

- `cmd/api`: Contains the main entry point of the application.
- `internal/api`: Contains the API handlers and router.
- `internal/config`: Contains the configuration loading logic.
- `internal/domain`: Contains the domain models.
- `internal/repository`: Contains the repository implementations for interacting with the database.
- `internal/services`: Contains the service layer implementations.
- `migrations`: Contains the database migration scripts.
- `scripts`: Contains helper scripts for building, running, and testing the application.

## Prerequisites

- Docker
- Docker Compose

## Running the Application

### Using Docker Compose

1. **Start the postgres containers:**

   ```sh
   docker-compose up --build
   ```

   This command will build the Docker image and start the container for the PostgreSQL database.

2. **Build the appluication:**

   ```sh
   docker run --rm -v $(pwd):/mnt -w /mnt order-service-test /mnt/scripts/build.sh
   ```

   This script will build the Go application.

3. **Test the application:**

   ```sh
   docker run --rm -v $(pwd):/mnt -w /mnt order-service-test /mnt/scripts/test.sh
   ```

   This script will run the tests and generate a coverage report.

4. **Apply database migrations and run the application:**

   ```sh
   docker run --rm -v $(pwd):/mnt \
   -e DB_HOST=host.docker.internal \
   -e DB_PORT=5432 \
   -e DB_USER=postgres \
   -e DB_PASSWORD=postgres \
   -e DB_NAME=order_service \
   -p 9090:9090 \
   -w /mnt order-service-test /mnt/scripts/run.sh
   ```

   The `run.sh` script in the `scripts` directory will automatically apply the database migrations when the application container starts.

5. **Access the application:**

   The application will be accessible at `http://localhost:9090`.

### Running Scripts

- **Build the application:**

  ```sh
  ./scripts/build.sh
  ```

  This script will build the Go application.

- **Run the application:**

  ```sh
  ./scripts/run.sh
  ```

  This script will apply the database migrations and start the application.

- **Run tests:**

  ```sh
  ./scripts/test.sh
  ```

  This script will run the tests and generate a coverage report.

## API Endpoints

- **Create Order:**

  ```
  POST /api/orders
  ```

  Request body example:

  ```json
  {
    "order": {
      "items": [
        { "product_id": 1, "quantity": 2 },
        { "product_id": 2, "quantity": 3 }
      ]
    }
  }
  ```

- **Get Order:**

  ```
  GET /api/orders/{id}
  ```

  Response example:

  ```json
  {
    "order_id": 1,
    "order_price": 35.0,
    "order_vat": 3.5,
    "items": [
      { "product_id": 1, "quantity": 2, "price": 20.0, "vat": 2.0 },
      { "product_id": 2, "quantity": 3, "price": 15.0, "vat": 1.5 }
    ]
  }
  ```

## Configuration

The application configuration is loaded from environment variables. The following variables can be set:

- `PORT`: The port on which the application will run (default: `9090`).
- `DB_HOST`: The database host (default: `localhost`).
- `DB_PORT`: The database port (default: `5432`).
- `DB_USER`: The database user (default: `postgres`).
- `DB_PASSWORD`: The database password (default: `postgres`).
- `DB_NAME`: The database name (default: `order_service`).
- `DB_SSLMODE`: The database SSL mode (default: `disable`).
- `ENV`: The application environment (default: `development`).
- `LOG_LEVEL`: The log level (default: `info`).
