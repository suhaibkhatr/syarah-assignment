# Gift Store Service

A Go-based microservice for managing gift store operations with MySQL, Kafka, and Elasticsearch integration.

## Prerequisites

- Docker and Docker Compose
- Go 1.16+ (for local development)

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd gift-store
   ```

2. Create a `config.yaml` file based on the example configuration:
   ```yaml
   # Copy and customize this configuration
   # Make sure to update credentials and connection details
   ```

3. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

4. Wait for all services to be fully started before proceeding.

## Services

- **MySQL**: Database service running on port 3306
- **Kafka**: Message broker running on port 9092
- **Kafka Connect**: Debezium connector service on port 8083
- **Elasticsearch**: Search and analytics engine on port 9200
- **Gift Store**: Main application service

## Configuration

Create a `config.yaml` file in the root directory with the following structure:

```yaml
database:
  host: mysql
  port: 3306
  name: giftdb
  user: root
  password: rootpass

kafka:
  brokers: kafka:9092
  topics:
    - gift_events

elasticsearch:
  url: http://elasticsearch:9200
  index: gifts
```

## Development

### Building the Application

```bash
go build -o gift-store cmd/app/main.go
```

### Running Tests

```bash
go test ./...
```

## Docker

The application is containerized using Docker. The `docker-compose.yml` file defines all the necessary services.

### Start Services

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### View Logs

```bash
docker-compose logs -f
```

## License

[Specify License]
