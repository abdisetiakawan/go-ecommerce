
# Go E-commerce Project

## Overview

This project is a Go-based e-commerce application that follows the principles of clean architecture. It is designed to be maintainable, testable, and scalable, separating concerns into distinct layers. The application uses MySQL as its database for reliable data storage and management. Additionally, it incorporates Kafka for asynchronous processing related to orders, enhancing the application's ability to handle high volumes of transactions efficiently.

## Clean Architecture

The project is structured into several layers:

* **Delivery** : Contains HTTP handlers that manage incoming requests and responses.
* **Use Cases** : Contains business logic and application rules.
* **Repositories** : Handles data access and persistence.
* **Entities** : Represents the core data models of the application.

### Advantages of Clean Architecture

* **Separation of Concerns** : Each layer has a specific responsibility, making the codebase easier to navigate and maintain.
* **Testability** : Business logic is separated from the delivery mechanism, allowing for easier unit testing.
* **Flexibility** : Changes in one layer (e.g., switching from HTTP to gRPC) can be made with minimal impact on other layers.

## Setting Up the Project

### 1. Configure Environment Variables

Copy the example environment file and modify it as needed:

```sh
cp .env.example .env
```

### 2. Running the Project with Docker Compose

To start the project using Docker Compose, run:

```sh
docker-compose up -d
```

To start the development environment with hot reload and debugging enabled, use:

```sh
docker-compose -f docker-compose.dev.yaml up -d
```

### 3. Running Database Migrations

After the containers are up and running, execute the migration inside the running application container:

```sh
docker exec go-ecommerce make migrate
```

## Available Endpoints

### Authentication

* `POST /api/auth/register`: Register a new user.
* `POST /api/auth/login`: Log in an existing user.

### User Management

* `POST /api/user/profile`: Create a user profile.
* `GET /api/user/profile`: Retrieve user profile.
* `PUT /api/user/profile`: Update user profile.
* `PATCH /api/user/password`: Change user password.

### Buyer Operations

* `GET /api/buyer/orders`: Search for orders.
* `GET /api/buyer/orders/:order_uuid`: Get order details.
* `POST /api/buyer/orders`: Create a new order.
* `PATCH /api/buyer/orders/:order_uuid/cancel`: Cancel an order.
* `PATCH /api/buyer/orders/:order_uuid/checkout`: Checkout an order.

### Seller Operations

* `POST /api/seller/store`: Register a new store.
* `GET /api/seller/store`: Retrieve store details.
* `PUT /api/seller/store`: Update store information.
* `POST /api/seller/products`: Register a new product.
* `GET /api/seller/products`: Retrieve products.
* `GET /api/seller/products/:product_uuid`: Get product details.
* `PUT /api/seller/products/:product_uuid`: Update product information.
* `DELETE /api/seller/products/:product_uuid`: Delete a product.
* `GET /api/seller/orders`: Retrieve seller's orders.
* `GET /api/seller/orders/:order_uuid`: Get order details for seller.
* `PATCH /api/seller/orders/:order_uuid/shipping`: Update shipping status.

## Kafka Integration

The application utilizes Kafka for asynchronous processing of order-related tasks. This allows for better scalability and responsiveness, as order processing can occur in the background without blocking the main application flow.

## Swagger Documentation

The application includes Swagger documentation that can be accessed at the following URL: `/swagger`. This provides an interactive interface for testing the API endpoints and understanding the request/response structures.

## Future Enhancements

This project is currently in active development and will continue to evolve with the following planned enhancements:

1. **Redis Caching** : Improve response times and reduce database load by implementing caching for frequently accessed data.
2. **Scalability Improvements** : Add support for horizontal scaling and load balancing to handle increasing user traffic.
3. **Additional Features** : Include support for advanced search filters, analytics, and more.

## Conclusion

This documentation provides an overview of the Go E-commerce project, its architecture, and the available endpoints. The clean architecture approach ensures that the application is well-structured and easy to maintain. The integration of Kafka for asynchronous order processing and the availability of Swagger documentation further enhance the usability and performance of the application. As this project is still under development, contributions and suggestions are welcome to help it reach its full potential.
