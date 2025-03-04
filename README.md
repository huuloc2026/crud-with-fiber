The API will be available at http://localhost:3000 with the following endpoints:
Users:
POST /api/users - Create user
GET /api/users - Get all users
GET /api/users/:id - Get specific user
PUT /api/users/:id - Update user
DELETE /api/users/:id - Delete user
Products:
POST /api/products - Create product
GET /api/products - Get all products
GET /api/products/:id - Get specific product
PUT /api/products/:id - Update product
DELETE /api/products/:id - Delete product
You can test the API using tools like Postman or curl. For example, to create a new user:
This is a basic setup that you can build upon. You might want to add:

1. Input validation
   Authentication/Authorization
   Error handling middleware
   Password hashing
2. Environment configuration
   Logging
   Testing
