# PT XYZ Multifinance API

A comprehensive REST API for multifinance services built with Go, providing customer management, credit limit management, and transaction processing capabilities.

## Features

- **User Authentication**: JWT-based authentication with role-based access control (Admin/Customer)
- **Customer Management**: Complete customer profile management with KTP and selfie photo support
- **Credit Limit System**: Flexible tenor-based credit limits (1, 2, 3, 4 months)
- **Transaction Processing**: End-to-end transaction management with real-time limit validation
- **Security**: Rate limiting, CORS protection, and comprehensive input validation
- **Database**: MySQL with GORM ORM and automatic migrations

## Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin HTTP Framework
- **Database**: MySQL 8.0
- **ORM**: GORM
- **Authentication**: JWT with bcrypt password hashing
- **Containerization**: Docker with multi-stage builds
- **Testing**: Testify with mock repositories

## Quick Start

### Prerequisites

- Go 1.22 or higher
- MySQL 8.0
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/your-org/multifinance-xyz.git
cd multifinance-xyz
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Configure your environment variables in `.env`:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=kredit_plus

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=24
```

**For Docker usage**: If using Docker with host network, keep `DB_HOST=localhost`. If using bridge network, change to `DB_HOST=host.docker.internal` (Mac/Windows) or your host IP address (Linux).

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Using Docker

**Prerequisites**: Ensure MySQL is running on your host system:
```bash
sudo systemctl start mysql
sudo systemctl enable mysql
```

1. **Build and run with Docker Compose** (Recommended):
```bash
docker-compose up --build
```

2. **Or build the Docker image manually**:
```bash
# Build the image
docker build -t multifinance-api .

# Run with host network (for MySQL access)
docker run -d \
  --name multifinance-app \
  --network host \
  --env-file .env \
  -e DB_HOST=localhost \
  multifinance-api

# Check logs
docker logs multifinance-app

# Test the application
curl http://localhost:8080/health
```

3. **Alternative: Run without host network**:
```bash
# For systems that support host.docker.internal
docker run -d \
  --name multifinance-app \
  -p 8080:8080 \
  --env-file .env \
  -e DB_HOST=host.docker.internal \
  multifinance-api

# For Linux systems, use host IP
docker run -d \
  --name multifinance-app \
  -p 8080:8080 \
  --env-file .env \
  -e DB_HOST=172.17.0.1 \
  multifinance-api
```

**Note**: The application connects to MySQL running on the host system, not in a container. Make sure your MySQL service is accessible and the `kredit_plus` database exists.

## API Documentation

### Authentication

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Base URL
```
http://localhost:8080/api/v1
```

### Public Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123",
  "confirm_password": "SecurePass123",
  "role": "CUSTOMER",
  "customer_data": {
    "nik": "3174012345678901",
    "full_name": "John Doe",
    "legal_name": "John Doe",
    "birth_place": "Jakarta",
    "birth_date": "1990-01-15",
    "salary": 8000000,
    "ktp_photo_path": "/uploads/ktp/john_ktp.jpg",
    "selfie_photo_path": "/uploads/selfie/john_selfie.jpg",
    "limits": [
      {"tenor_months": 1, "limit_amount": 1000000},
      {"tenor_months": 2, "limit_amount": 2000000},
      {"tenor_months": 3, "limit_amount": 3000000},
      {"tenor_months": 4, "limit_amount": 4000000}
    ]
  }
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePass123"
}
```

### Customer Endpoints

#### Get My Profile
```http
GET /customers/me
Authorization: Bearer <token>
```

#### Get Customer Limits
```http
GET /customers/{id}/limits
Authorization: Bearer <token>
```

### Transaction Endpoints

#### Create Transaction
```http
POST /transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "customer_id": 1,
  "tenor_months": 1,
  "otr_amount": 500000,
  "admin_fee": 50000,
  "interest_amount": 25000,
  "asset_name": "iPhone 15 Pro",
  "asset_type": "WHITE_GOODS",
  "transaction_source": "ECOMMERCE"
}
```

#### Get Customer Transactions
```http
GET /transactions/customer/{customer_id}
Authorization: Bearer <token>
```

### Admin Endpoints

#### Get All Customers
```http
GET /admin/customers?limit=10&offset=0
Authorization: Bearer <admin-token>
```

#### Get All Transactions
```http
GET /admin/transactions?limit=10&offset=0
Authorization: Bearer <admin-token>
```

#### Update Transaction Status
```http
PUT /admin/transactions/{id}/status
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "status": "APPROVED"
}
```

## Business Rules

### Customer Registration
- Each customer must provide exactly 4 tenor limits (1, 2, 3, 4 months)
- NIK must be unique and exactly 16 digits
- All required fields must be provided

### Transaction Processing
- Transactions are created with PENDING status
- Amount must not exceed available credit limit for the specified tenor
- Contract numbers are auto-generated with format: XYZ{timestamp}
- Installment amount is automatically calculated

### Credit Limits
- PT XYZ only supports 1, 2, 3, and 4-month tenors
- Each customer must have limits for all four tenors
- Used amounts are updated in real-time during transaction creation
- Limits are rolled back if transactions are rejected

## Database Schema

### Users Table
- User authentication and role management
- Supports ADMIN and CUSTOMER roles

### Customers Table
- Customer profile information
- Links to Users table via foreign key
- Stores KTP and selfie photo paths

### Customer Limits Table
- Credit limits per tenor for each customer
- Tracks used and available amounts
- Unique constraint on customer_id + tenor_months

### Transactions Table
- Complete transaction records
- Links to customers and tracks limit usage
- Supports multiple asset types and transaction sources

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./test/usecase/
```

### Code Structure
```
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain entities and interfaces
│   ├── infrastructure/  # Database and external services
│   ├── interfaces/      # HTTP handlers and DTOs
│   └── usecase/         # Business logic
├── pkg/                 # Shared utilities
├── test/                # Test files and mocks
└── main.go             # Application entry point
```

### Default Admin Account
When the application starts, it automatically creates a default admin account:
- Username: `admin`
- Password: `admin123`
- Email: `admin@ptxyz.com`

**Important**: Change this password in production environments.

## Production Deployment

### Environment Variables
Ensure all required environment variables are set:
```env
SERVER_PORT=8080
GIN_MODE=release
DB_HOST=your-db-host
DB_PASSWORD=secure-password
JWT_SECRET=your-secure-secret
BCRYPT_COST=12
```

### Health Check
The application provides a health check endpoint:
```http
GET /health
```

### Security Considerations
- Use strong JWT secrets in production
- Enable HTTPS in production
- Configure proper CORS origins
- Set appropriate rate limits
- Use environment-specific database credentials

## Troubleshooting

### Docker Connection Issues

If you encounter database connection errors when running with Docker:

1. **"connection refused" errors**:
```bash
# Ensure MySQL is running on host
sudo systemctl status mysql

# Use host network mode
docker run -d --name multifinance-app --network host --env-file .env -e DB_HOST=localhost multifinance-api
```

2. **"no such host" errors**:
```bash
# For Linux systems, use host IP instead of host.docker.internal
docker run -d --name multifinance-app -p 8080:8080 --env-file .env -e DB_HOST=172.17.0.1 multifinance-api
```

3. **Check container logs**:
```bash
docker logs multifinance-app
docker logs -f multifinance-app  # Follow logs real-time
```

4. **Container management**:
```bash
# Stop and remove container
docker stop multifinance-app
docker rm multifinance-app

# Remove and rebuild image
docker rmi multifinance-api
docker build -t multifinance-api .
```

### Application Issues

For general application debugging:
- Check if the database `kredit_plus` exists
- Verify MySQL user permissions
- Ensure all required environment variables are set
- Check application logs for specific error messages

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is proprietary software of PT XYZ Multifinance.