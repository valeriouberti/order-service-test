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

## Design Considerations

This order service project was designed with a microservices mindset, focusing on modularity, maintainability. The service follows a clear separation of concerns across layers—routing, service, and repository—to facilitate ease of testing and future enhancements.

A key design decision was to use Go’s built-in net/http library instead of a more opinionated framework such as Gin. The reasoning behind this was twofold. First, the net/http package offers developers full visibility and control over the HTTP handling pipeline without hidden abstractions. Second, reducing external dependencies for simple project like this, ensuring that the service remains lightweight and straightforward to maintain.

In summary, this design approach was chosen to ensure clarity and efficiency while keeping the code flexible.

## Database Considerations

PostgreSQL was chosen as the database backend due to its robustness, reliability, and extensive feature set. Here are the key reasons for its selection:

- **Stability and Performance:** PostgreSQL offers a high level of data consistency and can efficiently handle complex queries and high transaction volumes, making it ideal for an order service that may scale over time.
- **ACID Compliance:** With full support for ACID transactions, PostgreSQL ensures that all order-related operations are executed reliably, guaranteeing data integrity even under concurrent access.
- **Extensibility and Features:** The rich ecosystem of data types, indexing options, and support for custom functions allows for flexible data manipulation and optimizations as the application evolves.
- **Community and Ecosystem:** Being widely adopted, PostgreSQL benefits from a large community, active development, and frequent updates, which provides a level of assurance and access to a vast array of extensions and tools.
- **Ease of Integration:** The integration with Go is straightforward via robust PostgreSQL drivers. This minimizes boilerplate code and leverages Docker Compose for easy local provisioning and testing.

## Future Improvements

- **Authentication Middleware:**  
  Implement an authentication middleware to secure API endpoints. This middleware will validate user tokens or api keys to restrict access based on authorization levels, ensuring that only authenticated users or services can perform sensitive operations.

- **Product Management API:**  
  Introduce new API endpoints to manage products, including inserting, updating, and deleting product records. This new functionality will enable dynamic updates to the product catalog.

- **Enhanced Error Handling:**  
  Standardize error responses using custom error objects. This will improve client-side error processing and streamline debugging.

- **API Documentation:**  
  Integrate tools like Swagger or OpenAPI to generate interactive API documentation, making it easier for developers and external integrators to understand and use the API.

- **Performance & Scalability Enhancements:**  
  Consider adding database indexing improvements and caching mechanisms to further boost performance as the service scales.

- **Monitoring and Logging:**  
  Incorporate advanced logging and monitoring solutions to gain deeper insights into application performance and errors, which can aid in proactive troubleshooting.
