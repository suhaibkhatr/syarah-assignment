# Gift Store Service

A Go-based microservice for managing gift store operations with MySQL, Kafka, and Elasticsearch integration. The system provides:

- Real-time data synchronization between MySQL and Elasticsearch using Debezium and Kafka
- RESTful API for managing products
- Scalable and distributed architecture

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

4. **Important**: After starting the services, you need to set up the Debezium MySQL connector. Run the following command to configure it:
  ```bash
  curl --location 'http://localhost:8083/connectors' \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "gift_store-connector",
    "config": {
      "connector.class": "io.debezium.connector.mysql.MySqlConnector",
      "database.hostname": "mysql",
      "database.port": "3306",
      "database.user": "root",
      "database.password": "rootpass",
      "database.server.id": "184054",
      "database.include.list": "gift_store",
      "topic.prefix": "gift_store",
      "schema.history.internal.kafka.bootstrap.servers": "kafka:9094",
      "schema.history.internal.kafka.topic": "schema-changes.gift_store"
    }
  }'
  ```
   This step is crucial for enabling real-time data synchronization between MySQL and Elasticsearch.

5. The application will start and connect to the required services.

## Services

- **MySQL**: Database service running on port 3306
- **Kafka**: Message broker running on port 9092
- **Kafka Connect**: Debezium connector service on port 8083
- **Elasticsearch**: Search and analytics engine on port 9200
- **Gift Store**: Main application service
- **Gift Store API**: REST API service running on port 8080

## Configuration

### Environment Variables

Configuration is managed through environment variables. Copy `.env.example` to `.env` and update the values as needed:

```env
# MySQL Configuration
MYSQL_DNS=root:rootpass@tcp(mysql:3306)/gift_store?parseTime=true

# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_TOPICS=gift_store.gift_store.products
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
  "name": "gift_store-connector",
  "config": {
    "connector.class": "io.debezium.connector.mysql.MySqlConnector",
    "database.hostname": "mysql",
    "database.port": "3306",
    "database.user": "root",
    "database.password": "rootpass",
    "database.server.id": "184054",
    "database.include.list": "gift_store",
    "topic.prefix": "gift_store",
    "schema.history.internal.kafka.bootstrap.servers": "kafka:9094",
    "schema.history.internal.kafka.topic": "schema-changes.gift_store"
  }
}'
```

This will create a Debezium MySQL connector that:
- Connects to MySQL at `mysql:3306`
- Monitors the `gift_store` database
- Writes changes to Kafka topics with the prefix `gift_store`
- Stores schema history in a Kafka topic

You can verify the connector is running by visiting:
- Connector status: http://localhost:8083/connectors/gift_store-connector/status
- List all connectors: http://localhost:8083/connectors

### Environment Variables

- `MYSQL_DNS`: MySQL connection string
- `KAFKA_BROKERS`: Comma-separated list of Kafka brokers
- `KAFKA_TOPICS`: Comma-separated list of Kafka topics to subscribe to
- `KAFKA_GROUP_ID`: Consumer group ID for Kafka
- `ELASTIC_URL`: Elasticsearch server URL

## API Endpoints

### Products

#### Get All Products
```
GET /api/v1/products
```

**Response**
```json
[
  {
    "id": 1,
    "name": "Gift Card",
    "description": "$50 Gift Card",
    "price": 50.0,
    "category": "Gift Cards",
    "gift_category": "General",
    "age_group": "All Ages",
    "brand": "Gift Store",
    "is_available": true,
    "created_at": "2025-07-07T00:00:00Z",
    "updated_at": "2025-07-07T00:00:00Z"
  }
]
```

#### Get Product by ID
```
GET /api/v1/products/:id
```

**Response**
```json
{
  "id": 1,
  "name": "Gift Card",
  "description": "$50 Gift Card",
  "price": 50.0,
  "category": "Gift Cards",
  "gift_category": "General",
  "age_group": "All Ages",
  "brand": "Gift Store",
  "is_available": true,
  "created_at": "2025-07-07T00:00:00Z",
  "updated_at": "2025-07-07T00:00:00Z"
}
```

#### Create Product
```
POST /api/v1/products
```

**Request Body**
```json
{
  "name": "New Gift Card",
  "description": "$100 Gift Card",
  "price": 100.0,
  "category": "Gift Cards",
  "gift_category": "Premium",
  "age_group": "All Ages",
  "brand": "Gift Store",
  "is_available": true
}
```

**Response**
```json
{
  "id": 2,
  "name": "New Gift Card",
  "description": "$100 Gift Card",
  "price": 100.0,
  "category": "Gift Cards",
  "gift_category": "Premium",
  "age_group": "All Ages",
  "brand": "Gift Store",
  "is_available": true,
  "created_at": "2025-07-07T00:00:00Z",
  "updated_at": "2025-07-07T00:00:00Z"
}
```

#### Update Product
```
PUT /api/v1/products/:id
```

**Request Body**
```json
{
  "name": "Updated Gift Card",
  "description": "$100 Gift Card (Updated)",
  "price": 99.99,
  "is_available": false
}
```

**Response**
```json
{
  "message": "Product updated successfully",
  "product": {
    "id": 2,
    "name": "Updated Gift Card",
    "description": "$100 Gift Card (Updated)",
    "price": 99.99,
    "category": "Gift Cards",
    "gift_category": "Premium",
    "age_group": "All Ages",
    "brand": "Gift Store",
    "is_available": false,
    "created_at": "2025-07-07T00:00:00Z",
    "updated_at": "2025-07-07T01:00:00Z"
  }
}
```

#### Delete Product
```
DELETE /api/v1/products/:id
```

**Response**
```json
{
  "message": "Product deleted successfully"
}
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
