<div align="center">
<img src="https://media2.dev.to/dynamic/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Fyzdi87bkixukz6kz5a0m.png"  width="100%">
</div>

# Fiber CRUD API

A modern Go web application built with Fiber framework, featuring:
- RESTful API endpoints
- JWT Authentication
- PostgreSQL database
- Redis caching
- RabbitMQ for async tasks
- Docker containerization
- CI/CD pipeline with GitHub Actions

## Features

- User authentication (register/login) with JWT
- Session caching with Redis
- Asynchronous task processing with RabbitMQ
- Containerized with Docker and Docker Compose
- Automated testing and deployment with GitHub Actions

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/fiber-crud.git
cd fiber-crud
```

2. Create a config file:
```bash
cp config/config.example.yaml config/config.yaml
```

3. Start the application with Docker Compose:
```bash
docker-compose up -d
```

The application will be available at `http://localhost:3000`

## Environment Variables

The following environment variables can be set in your config.yaml file:

```yaml
server:
  host: "0.0.0.0"
  port: 3000

database:
  host: "postgres"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "myapp"

redis:
  host: "redis"
  port: "6379"
  password: ""
  db: 0

rabbitmq:
  host: "rabbitmq"
  port: "5672"
  user: "guest"
  password: "guest"

jwt:
  secret: "your-secret-key"
  expiresIn: "24h"
```

## API Endpoints

### Authentication
- POST `/auth/register` - Register a new user
- POST `/auth/login` - Login and receive JWT token

### Protected Routes (Requires JWT Token)
- GET `/api/users` - Get all users
- POST `/api/users` - Create a new user
- GET `/api/users/:id` - Get user by ID
- PUT `/api/users/:id` - Update user
- DELETE `/api/users/:id` - Delete user

## Development

To run the application locally:

1. Install dependencies:
```bash
go mod download
```

2. Start the required services:
```bash
docker-compose up -d postgres redis rabbitmq
```

3. Run the application:
```bash
go run main.go
```

## Testing

Run the tests:
```bash
go test -v ./...
```

## CI/CD

The project includes a GitHub Actions workflow that:
1. Runs tests
2. Builds the application
3. Pushes the Docker image to Docker Hub (on main branch)

To enable CI/CD:
1. Fork the repository
2. Add the following secrets to your GitHub repository:
   - DOCKER_HUB_USERNAME
   - DOCKER_HUB_TOKEN

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
