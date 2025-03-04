<div align="center">
<img src="https://media2.dev.to/dynamic/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Fyzdi87bkixukz6kz5a0m.png"  width="100%">
</div>

## MyApp

MyApp is a RESTful API built with Go, Fiber, and PostgreSQL. It provides CRUD operations for users and products, with JWT-based authentication.

## Features

- User registration and login with JWT authentication
- CRUD operations for users and products
- Pagination support for listing users and products
- Password hashing for secure storage

## Prerequisites

- Go 1.16 or later
- PostgreSQL
- Git

## Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/huuloc2026/crud-with-fiber
   cd myapp
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Configure the database:**

   Update the `config/config.yaml` file with your PostgreSQL credentials:

   ```yaml
   database:
     host: localhost
     port: 5433
     user: root
     password: huuloc2026
     name: fiber-api
     sslmode: disable
     timezone: Asia/Jakarta

   jwt:
     secret: "super-secret-key-change-in-production"
     expires_in: 24h
   ```

4. **Run the application:**

   ```bash
   go run main.go
   ```

   The server will start on the configured port (default is 3000).

## API Endpoints

### Authentication

- **POST /api/auth/register**: Register a new user
- **POST /api/auth/login**: Login and receive a JWT token

### Users

- **POST /api/users**: Create a new user (protected)
- **GET /api/users**: Get a list of users (protected)
- **GET /api/users/:id**: Get a specific user by ID (protected)
- **PUT /api/users/:id**: Update a user by ID (protected)
- **DELETE /api/users/:id**: Delete a user by ID (protected)

### Products

- **POST /api/products**: Create a new product (protected)
- **GET /api/products**: Get a list of products (protected)
- **GET /api/products/:id**: Get a specific product by ID (protected)
- **PUT /api/products/:id**: Update a product by ID (protected)
- **DELETE /api/products/:id**: Delete a product by ID (protected)

## Usage

1. **Register a new user:**

   ```bash
   curl -X POST http://localhost:3000/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"name":"John Doe","email":"john@example.com","password":"secret123"}'
   ```

2. **Login to get a JWT token:**

   ```bash
   curl -X POST http://localhost:3000/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"john@example.com","password":"secret123"}'
   ```

3. **Use the token for protected routes:**

   ```bash
   curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:3000/api/users
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
