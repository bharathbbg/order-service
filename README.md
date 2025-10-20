# Order Service

Part of the e-commerce platform microservices architecture. This service is responsible for order management.

## Features

- Create, read, update and delete orders
- Process payments
- Integrate with delivery and inventory services
- REST API for front-end communication
- gRPC API for inter-service communication

## Technology Stack

- Go with Gin framework for REST API
- gRPC for inter-service communication
- PostgreSQL for data persistence
- Redis for caching
- Docker for containerization

## Local Development

```bash
# Start dependencies
docker-compose up -d postgres redis

# Run service
go run cmd/server/main.go
```

## API Documentation

### REST Endpoints

- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - List all orders
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id` - Update an order
- `DELETE /api/v1/orders/:id` - Delete an order

### gRPC Services

See `proto/order.proto` for full service definition.