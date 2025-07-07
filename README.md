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

2. Copy the example environment file and update the values:
   ```bash
   cp .env.example .env
   ```
   Edit the `.env` file with your configuration.

3. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

4. The application will start and connect to the required services.

## Services

- **MySQL**: Database service running on port 3306
- **Kafka**: Message broker running on port 9092
- **Kafka Connect**: Debezium connector service on port 8083
- **Elasticsearch**: Search and analytics engine on port 9200
- **Gift Store**: Main application service

## Configuration

### Environment Variables

Configuration is managed through environment variables. Copy `.env.example` to `.env` and update the values as needed:

```env
# MySQL Configuration
MYSQL_DNS=root:rootpass@tcp(mysql:3306)/giftdb?parseTime=true

# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPICS=giftdb.giftdb.products
KAFKA_GROUP_ID=gift-store-group

# Elasticsearch Configuration
ELASTIC_URL=http://elasticsearch:9200
```

### Debezium MySQL Connector Setup

To set up the Debezium MySQL connector, run the following command after starting the services:

```bash
curl --location 'http://localhost:8083/connectors' \
--header 'Content-Type: application/json' \
--data '{
  "name": "giftdb-connector",
  "config": {
    "connector.class": "io.debezium.connector.mysql.MySqlConnector",
    "database.hostname": "mysql",
    "database.port": "3306",
    "database.user": "root",
    "database.password": "rootpass",
    "database.server.id": "184054",
    "database.include.list": "giftdb",
    "topic.prefix": "giftdb",
    "schema.history.internal.kafka.bootstrap.servers": "kafka:9094",
    "schema.history.internal.kafka.topic": "schema-changes.giftdb"
  }
}'
```

This will create a Debezium MySQL connector that:
- Connects to MySQL at `mysql:3306`
- Monitors the `giftdb` database
- Writes changes to Kafka topics with the prefix `giftdb`
- Stores schema history in a Kafka topic

You can verify the connector is running by visiting:
- Connector status: http://localhost:8083/connectors/giftdb-connector/status
- List all connectors: http://localhost:8083/connectors

### Environment Variables

- `MYSQL_DNS`: MySQL connection string
- `KAFKA_BROKERS`: Comma-separated list of Kafka brokers
- `KAFKA_TOPICS`: Comma-separated list of Kafka topics to subscribe to
- `KAFKA_GROUP_ID`: Consumer group ID for Kafka
- `ELASTIC_URL`: Elasticsearch server URL

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
