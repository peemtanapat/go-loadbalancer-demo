# Go Load Balancer Demo

A simple implementation of a custom load balancer in Go using round-robin algorithm. This project demonstrates how to build a load balancer from scratch that distributes incoming requests across multiple backend API services.

## Features

- **Round-robin load balancing** - Distributes requests evenly across available servers
- **Health checking** - Monitors backend server health and excludes unhealthy servers
- **Dockerized setup** - Easy deployment with Docker Compose
- **REST API endpoints** - Includes sample API services for testing
- **Load balancer status** - Real-time monitoring of load balancer state

## Prerequisites

Before running this project, make sure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) (version 20.0 or higher)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.0 or higher)
- [Go](https://golang.org/doc/install) (version 1.19 or higher) - for development

## Project Structure

```
.
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile.apiservice       # Dockerfile for API services
├── Dockerfile.loadbalancer     # Dockerfile for load balancer
├── .env                        # Environment variables configuration
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies
├── rest.http                   # HTTP requests for testing
├── api/
│   └── user_api.go            # Sample API service implementation
└── loadbalancer/
    └── loadbalancer.go        # Load balancer implementation
```

## Quick Start

### 1. Clone and Navigate to Project

```bash
git clone <repository-url>
cd go-loadbalancer-demo
```

### 2. Configure Environment (Optional)

The project uses a `.env` file to configure the load balancer's target services. The default configuration is:

```bash
# Load Balancer Configuration
TARGET_SERVICES=http://host.docker.internal:8081,http://host.docker.internal:8082,http://host.docker.internal:8083
```

You can modify the `.env` file to change the target services or create environment-specific configurations:

```bash
# For production
cp .env .env.prod
# Edit .env.prod with production service URLs

# For development
cp .env .env.dev
# Edit .env.dev with development service URLs
```

### 3. Start the Services

```bash
docker-compose up --build
```

This command will:

- Load configuration from `.env` file
- Build and start 3 API service instances (ports 8081, 8082, 8083)
- Build and start the load balancer (port 9080)
- Create a Docker network for inter-service communication

### 4. Verify the Setup

Open your browser or use curl to check the load balancer status:

```bash
curl http://localhost:9080/lb-status
```

### 5. Stop the Services

```bash
docker-compose down
```

## API Endpoints

The load balancer exposes the following endpoints:

### Load Balancer Status

- **GET** `http://localhost:9080/lb-status`
- Returns the current status of the load balancer and all backend servers

### API Endpoints (proxied through load balancer)

- **GET** `http://localhost:9080/api/users` - Get list of users
- **POST** `http://localhost:9080/api/users` - Create a new user
- **GET** `http://localhost:9080/api/heavy-task` - Simulate a heavy processing task

## Example Requests and Responses

### 1. Check Load Balancer Status

**Request:**

```http
GET http://localhost:9080/lb-status
```

**Response:**

```json
{
  "loadBalancer": "active",
  "servers": [
    {
      "url": {
        "Scheme": "http",
        "Opaque": "",
        "User": null,
        "Host": "host.docker.internal:8081",
        "Path": "",
        "RawPath": "",
        "OmitHost": false,
        "ForceQuery": false,
        "RawQuery": "",
        "Fragment": "",
        "RawFragment": ""
      },
      "healthy": true
    },
    {
      "url": {
        "Scheme": "http",
        "Opaque": "",
        "User": null,
        "Host": "host.docker.internal:8082",
        "Path": "",
        "RawPath": "",
        "OmitHost": false,
        "ForceQuery": false,
        "RawQuery": "",
        "Fragment": "",
        "RawFragment": ""
      },
      "healthy": true
    },
    {
      "url": {
        "Scheme": "http",
        "Opaque": "",
        "User": null,
        "Host": "host.docker.internal:8083",
        "Path": "",
        "RawPath": "",
        "OmitHost": false,
        "ForceQuery": false,
        "RawQuery": "",
        "Fragment": "",
        "RawFragment": ""
      },
      "healthy": true
    }
  ],
  "algorithm": "round-robin",
  "timestamp": "2025-09-06T11:23:57.905241803Z"
}
```

### 2. Get Users

**Request:**

```http
GET http://localhost:9080/api/users
```

**Response:**

```json
{
  "port": "8080",
  "timestamp": "2025-09-06T11:28:37.631095502Z",
  "users": ["Alice", "Bird", "Charlie", "Dan"],
  "servedBy": "api-service-2"
}
```

### 3. Heavy Task Processing

**Request:**

```http
GET http://localhost:9080/api/heavy-task
```

**Response:**

```json
{
  "port": "8080",
  "timestamp": "2025-09-06T11:28:52.763918509Z",
  "servedBy": "api-service-3",
  "message": "Heavy task completed",
  "processingTimeMs": 2002
}
```

### 4. Create User

**Request:**

```http
POST http://localhost:9080/api/users
Content-Type: application/json

{
    "name": "Tanapat"
}
```

**Response:**

```json
{
  "port": "8080",
  "timestamp": "2025-09-06T11:29:05.507822709Z",
  "servedBy": "api-service-1",
  "message": "User created successfully",
  "user": {
    "name": "Tanapat"
  }
}
```

## How It Works

1. **Load Balancer**: Receives incoming requests on port 9080
2. **Round-Robin Algorithm**: Distributes requests sequentially across available backend servers
3. **Health Checking**: Periodically checks backend server health via `/health` endpoint
4. **Request Proxying**: Forwards requests to healthy backend servers and returns responses
5. **Automatic Failover**: Excludes unhealthy servers from the rotation

## Testing the Load Balancer

You can test the round-robin behavior by making multiple requests to the same endpoint and observing the `servedBy` field in the responses, which will rotate between `api-service-1`, `api-service-2`, and `api-service-3`.

## Development

To modify the code and rebuild:

1. Make your changes to the Go files
2. Rebuild and restart the services:
   ```bash
   docker-compose down
   docker-compose up --build
   ```

## Learning Points

This project demonstrates:

- Building a reverse proxy in Go
- Implementing round-robin load balancing
- Health checking and automatic failover
- Docker containerization and networking
- Inter-service communication
- HTTP request proxying and response handling

## Port Configuration

- **Load Balancer**: 9080
- **API Service 1**: 8081 (internal: 8080)
- **API Service 2**: 8082 (internal: 8080)
- **API Service 3**: 8083 (internal: 8080)

## Environment Variables

The project uses environment variables for configuration, managed through the `.env` file:

### Load Balancer (via .env file)

- `TARGET_SERVICES`: Comma-separated list of backend service URLs
  - Default: `http://host.docker.internal:8081,http://host.docker.internal:8082,http://host.docker.internal:8083`
  - You can modify this in the `.env` file to add/remove target services

### API Services (via docker-compose.yml)

- `PORT`: Service port (default: 8080)
- `INSTANCE_NAME`: Unique instance identifier

## Configuration Management

### Modifying Target Services

To change the backend services that the load balancer targets:

1. Edit the `.env` file:

   ```bash
   TARGET_SERVICES=http://host.docker.internal:8081,http://host.docker.internal:8082
   ```

2. Restart the services:
   ```bash
   docker-compose down
   docker-compose up
   ```

### Environment-Specific Configurations

For different environments, you can create separate env files:

```bash
# Development
echo "TARGET_SERVICES=http://localhost:8081,http://localhost:8082" > .env.dev

# Production
echo "TARGET_SERVICES=http://prod-api-1:8080,http://prod-api-2:8080,http://prod-api-3:8080" > .env.prod

# Use specific env file
docker-compose --env-file .env.prod up
```

## License

This project is for educational purposes and demonstrates load balancer implementation concepts.
