# Zion API

![Zion API Logo](https://via.placeholder.com/150?text=Zion+API)

## 🚀 Overview

Zion API is a modern, high-performance RESTful API built with Go and PostgreSQL. It follows Clean Architecture principles to ensure maintainability, testability, and scalability. The API provides robust endpoints for managing products and customers with authentication powered by Clerk.

## ✨ Features

- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **RESTful API**: Well-designed endpoints following REST principles
- **Authentication**: Secure JWT-based authentication with Clerk
- **Database**: PostgreSQL for reliable data persistence
- **Query Builder**: Type-safe SQL queries with Squirrel
- **Migration Tools**: Database versioning and migrations
- **Data Import**: Tools for importing legacy data from MySQL
- **CORS Support**: Cross-Origin Resource Sharing enabled
- **Pagination**: Efficient data retrieval with pagination support
- **Error Handling**: Comprehensive error handling and reporting

## 🏗️ Architecture

The project follows Clean Architecture principles with the following layers:

```
zion-api/
├── cmd/                  # Application entry points
│   ├── api/              # Main API server
│   └── scripts/          # Utility scripts (data migration, etc.)
├── internal/             # Private application code
│   ├── domain/           # Business entities and interfaces
│   ├── usecase/          # Business logic
│   ├── controller/       # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   └── infra/            # Infrastructure implementations
│       ├── database/     # Database connection and migrations
│       └── repository/   # Data access implementations
└── Makefile              # Build and development commands
```

## 🛠️ Technology Stack

- **Go**: Fast and efficient programming language
- **Gin**: High-performance HTTP web framework
- **PostgreSQL**: Advanced open-source relational database
- **pgx**: PostgreSQL driver and toolkit
- **Squirrel**: Fluent SQL generation library
- **JWT**: JSON Web Tokens for authentication
- **Clerk**: Authentication and user management
- **godotenv**: Environment variable management
- **testify**: Testing toolkit

## 🚦 Getting Started

### Prerequisites

- Go 1.18 or higher
- PostgreSQL 13 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/deividr/zion-api.git
   cd zion-api
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Run database migrations:

   ```bash
   make migration_up
   ```

5. Start the server:

   ```bash
   make run
   ```

The API will be available at `http://localhost:8000`.

## 📚 API Documentation

### Authentication

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### Endpoints

#### Products

| Method | Endpoint        | Description                         |
| ------ | --------------- | ----------------------------------- |
| GET    | `/products`     | List all products (with pagination) |
| GET    | `/products/:id` | Get product by ID                   |
| POST   | `/products`     | Create a new product                |
| PUT    | `/products/:id` | Update a product                    |
| DELETE | `/products/:id` | Delete a product (soft delete)      |

#### Customers

| Method | Endpoint         | Description                          |
| ------ | ---------------- | ------------------------------------ |
| GET    | `/customers`     | List all customers (with pagination) |
| GET    | `/customers/:id` | Get customer by ID                   |
| POST   | `/customers`     | Create a new customer                |
| PUT    | `/customers/:id` | Update a customer                    |
| DELETE | `/customers/:id` | Delete a customer (soft delete)      |

## 🧪 Testing

Run the test suite:

```bash
make test
```

Run tests with coverage:

```bash
make test_coverage
```

## 🛠️ Development

### Makefile Commands

| Command               | Description                           |
| --------------------- | ------------------------------------- |
| `make run`            | Start the API server                  |
| `make build`          | Build the application                 |
| `make test`           | Run tests                             |
| `make migration_up`   | Apply database migrations             |
| `make migration_down` | Rollback database migrations          |
| `make load_products`  | Import products from legacy database  |
| `make load_customers` | Import customers from legacy database |

### Adding a New Entity

1. Define the entity in the domain layer
2. Create repository interfaces in the domain layer
3. Implement the repository in the infrastructure layer
4. Create use cases in the usecase layer
5. Create controllers in the controller layer
6. Register routes in the main.go file

## 🔄 Data Migration

The project includes scripts to migrate data from a legacy MySQL database to PostgreSQL:

```bash
# Import products
make load_products

# Import customers
make load_customers
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👏 Acknowledgements

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [pgx - PostgreSQL Driver](https://github.com/jackc/pgx)
- [Squirrel - SQL Query Builder](https://github.com/Masterminds/squirrel)
- [Clerk - Authentication Provider](https://clerk.dev/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 👏 TODO

[] - Mover as entidades que estão soltas na pasta domain para um pasta entities;

---

Built with ❤️ by [Deivid Rodrigues](https://github.com/deividr)
