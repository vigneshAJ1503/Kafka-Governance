# Kafka Governance Control Plane

A stateless control-plane service for managing Kafka topic metadata and enforcing policy-based access control. This service acts as a governance layer between clients and Kafka clusters, enabling centralized topic management and fine-grained authorization using AWS Cedar.

## What It Does

- Stores and manages Kafka topic metadata (name, partitions, replication, ownership)
- Enforces policy-based access control using AWS Cedar policies
- Provides REST API for topic CRUD operations with authorization checks
- Acts as a control plane, not a data plane (does not produce/consume Kafka messages)
- Horizontally scalable and stateless

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────────────────────────┐
│         Routes (chi router)         │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│      API Layer (Handlers)           │
│   - topicApi.go                     │
│   - policyApi.go                    │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│    Service Layer (Business Logic)   │
│   - topic.go                        │
│   - policy.go                       │
└──────┬──────────────────────────────┘
       │
       ├─────────────────┬─────────────┐
       ▼                 ▼             ▼
┌──────────┐      ┌────────────┐ ┌─────────┐
│ MongoDB  │      │ Cedar CLI  │ │ Logger  │
│   (DB)   │      │ (Docker)   │ │ (Utils) │
└──────────┘      └────────────┘ └─────────┘
```

## Tech Stack

- **Language**: Go
- **HTTP Router**: chi
- **Database**: MongoDB
- **Policy Engine**: AWS Cedar (CLI via Docker)
- **Configuration**: Environment variables
- **Containerization**: Docker

## Project Structure

```
.
├── main.go               # Application entrypoint
├── go.mod                # Go module dependencies
├── docker-compose.yaml   # Docker compose configuration
├── Dockerfile            # Container definition
├── api/                  # HTTP request handlers
│   ├── topicApi.go
│   └── policyApi.go
├── service/              # Business logic layer
│   ├── topic.go
│   └── policy.go
├── db/                   # Database access layer
│   └── db.go
├── routes/               # HTTP route definitions
│   └── routes.go
├── models/               # Shared data structures
│   └── structs.go
├── config/               # Configuration loader
│   └── config.go
└── utils/                # Utilities
    ├── cedar.go          # Cedar policy integration
    ├── logger.go         # Logging
    └── response.go       # HTTP response helpers
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MongoDB (or use Docker Compose to run it)

### Setup

1. Clone the repository
```bash
git clone <repository-url>
cd Kafka-Governance
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables (see Configuration section)

4. Start MongoDB and Cedar CLI container
```bash
docker-compose up -d
```

### Run

#### Local Development
```bash
go run main.go
```

#### Using Docker
```bash
docker build -t kafka-governance .
docker run -p 8080:8080 --env-file .env kafka-governance
```

## Configuration

Configure the service using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DB_NAME` | Database name | `kafka_governance` |
| `PORT` | HTTP server port | `8080` |
| `CEDAR_CLI_ENDPOINT` | Cedar CLI Docker endpoint | `http://localhost:8180` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |

Example `.env` file:
```bash
MONGO_URI=mongodb://localhost:27017
DB_NAME=kafka_governance
PORT=8080
CEDAR_CLI_ENDPOINT=http://cedar-cli:8180
LOG_LEVEL=info
```

## API Endpoints

### Topics
- `POST /topics` - Create a new topic (with policy check)
- `GET /topics` - List all topics
- `GET /topics/{id}` - Get topic by ID
- `PUT /topics/{id}` - Update topic metadata
- `DELETE /topics/{id}` - Delete topic (with policy check)

### Policies
- `POST /policies` - Create a policy
- `GET /policies` - List all policies
- `DELETE /policies/{id}` - Delete a policy

## Scope & Notes

- **Control plane only**: This service manages topic metadata and enforces policies. It does not interact with Kafka brokers for message production/consumption.
- **Stateless**: All state is stored in MongoDB. No in-memory caching.
- **Horizontally scalable**: Multiple instances can run concurrently behind a load balancer.
- **Cedar integration**: Policy decisions are delegated to AWS Cedar CLI running in a Docker container.
- **No Kafka Admin API**: Topic creation/deletion in actual Kafka clusters must be handled separately (this service only manages metadata and authorization).

## Development

Run tests:
```bash
go test ./...
```

Build:
```bash
go build -o kafka-governance .
```

