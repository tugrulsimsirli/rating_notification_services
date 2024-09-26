
# Service Provider Ratings and Notifications

## Project Overview

This project consists of two backend API services for a service marketplace that allows customers to rate service providers and notifies them when new ratings are added. The two main services are:

- **Rating Service**: Handles customer ratings for service providers.
- **Notification Service**: Notifies service providers when they receive new ratings.

## Services

### 1. Rating Service
- **Endpoints**:
  - Submit a rating for a service provider.
  - Get the average rating of a specific service provider.
- **Database**: Stores ratings in a PostgreSQL database.
- **RabbitMQ**: Sends new rating data to the Notification Service via RabbitMQ.
  
### 2. Notification Service
- **Endpoints**:
  - Poll for the latest notifications about ratings.
  - Returns only new notifications (does not resend old notifications).
- **RabbitMQ**: Consumes rating notifications from the Rating Service.

## Technologies Used
- **Go**: Used to implement both Rating and Notification services.
- **PostgreSQL**: Used as the database for storing ratings.
- **RabbitMQ**: Used for communication between the services.
- **Docker**: Services are containerized using Docker.
- **Testcontainers**: For integration tests with RabbitMQ and PostgreSQL.
- **Swagger (OpenAPI)**: API documentation for both services.

## Project Setup

### 1. Clone the Repository
```bash
git clone <repository_url>
```

### 2. Build and Run with Docker
To build and run the services using Docker, ensure Docker and Docker Compose are installed.

```bash
make build
```

This will spin up the following containers:
- `rating_service`: Runs on port 8080.
- `notification_service`: Runs on port 8081.
- `db`: PostgreSQL database for storing ratings.
- `rabbitmq`: RabbitMQ message broker for communication.

### 3. Database Migration
Database schema will be initialized using the `db/migrations.sql` file during the container setup.

### 4. API Documentation (Swagger)
Once the services are running, you can access the Swagger API documentation:
- Rating Service: `http://localhost:8080/swagger/`
- Notification Service: `http://localhost:8081/swagger/`

## Testing

### Unit Tests
To run unit tests:
```bash
make test
```

### Integration Tests with Testcontainers
Testcontainers are used for running integration tests with RabbitMQ and PostgreSQL.

## Logging
Currently, the project does not include structured logging. This can be improved by integrating a logging framework like `logrus` or `zap`.

## CI/CD
The project does not currently include a CI/CD pipeline. However, you can integrate a pipeline using GitHub Actions or a similar tool for automating tests and releases.

## Improvements
- **Authentication**: Although the current implementation does not require authentication, adding security measures such as JWT or OAuth for authentication would be beneficial in a real-world scenario. Additionally, encrypting messages passed through RabbitMQ could enhance the overall security of the system.
- **Retry Mechanism**: Add a retry mechanism for RabbitMQ communication and other service interactions in case of failures. This will enhance the resilience of the system by ensuring communication is retried in case of temporary issues.
- **Logging**: Add structured logging for better error tracking and debugging.
- **CI/CD**: Implement CI/CD pipeline for automated tests and releases.
- **Advanced Error Handling**: Implement more granular error handling mechanisms.
- **Database Indexing**: Add indexes to frequently queried database columns to optimize query performance. This will reduce query execution times, especially as the dataset grows.
- **Caching**: Implement a caching mechanism (e.g., Redis) for average rating queries and notifications. This will reduce the number of database calls and improve performance, especially under heavy load.
- **Deployment Optimization**: While Docker is being used to containerize the services, migrating to Kubernetes could improve scalability and performance under high traffic. This would allow for better orchestration and auto-scaling.
- **ORM**: Implementing ORM for more efficient and maintainable database interactions.

