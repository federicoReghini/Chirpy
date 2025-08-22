# Chirpy 🐦

A modern, lightweight social media API built with Go. Chirpy provides a clean REST API for creating a Twitter-like social platform with user authentication, posts (chirps), and social interactions.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()
[![Go Report Card](https://img.shields.io/badge/Go%20Report-A+-brightgreen.svg)]()

## ✨ Features

- **User Management**: Registration, authentication, and profile management
- **Chirps (Posts)**: Create, read, update, and delete posts
- **Authentication**: JWT-based authentication system
- **RESTful API**: Clean REST endpoints with proper HTTP status codes
- **Database Integration**: Persistent data storage
- **Validation**: Comprehensive input validation and sanitization

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/chirpy.git
cd chirpy

# Install dependencies
go mod download

# Build the application
go build -o chirpy

# Run the server
./chirpy
```

The server will start on `http://localhost:8080` by default.

### Using Docker

```bash
# Build and run with Docker
docker build -t chirpy .
docker run -p 8080:8080 chirpy
```

## 📖 API Documentation

### Authentication

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Core Endpoints

#### Users

```http
POST /api/users          # Register a new user
POST /api/login          # Login user
PUT /api/users           # Update user profile
GET /api/users/{id}      # Get user by ID
```

#### Chirps

```http
GET /api/chirps          # Get all chirps
POST /api/chirps         # Create a new chirp
GET /api/chirps/{id}     # Get chirp by ID
DELETE /api/chirps/{id}  # Delete chirp (author only)
```

#### Health Check

```http
GET /api/healthz         # Health check endpoint
```

### Example Usage

#### Register a User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

#### Create a Chirp

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "body": "Hello, Chirpy world! 🐦"
  }'
```

#### Get All Chirps

```bash
curl http://localhost:8080/api/chirps
```

## 🏗️ Project Structure

```
chirpy/
├── main.go                 # Application entry point
├── handlers/               # HTTP request handlers
├── models/                 # Data models and structures
├── auth/                   # Authentication logic
├── database/               # Database connection and queries
├── middleware/             # HTTP middleware
├── config/                 # Configuration management
├── utils/                  # Utility functions
├── static/                 # Static files (HTML, CSS, JS)
├── tests/                  # Test files
├── docs/                   # Documentation
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksum
├── Dockerfile              # Docker configuration
└── README.md               # This file
```

## ⚙️ Configuration

The application can be configured using environment variables:

```bash
# Server configuration
PORT=8080                    # Server port (default: 8080)
HOST=localhost              # Server host (default: localhost)

# Database configuration
DB_URL=./database.db        # Database file path

# Authentication
JWT_SECRET=your-secret-key   # JWT signing secret
TOKEN_EXPIRY=24h            # Token expiration time

# Logging
LOG_LEVEL=info              # Log level (debug, info, warn, error)
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test ./handlers -run TestCreateChirp
```

## 🔧 Development

### Running in Development Mode

```bash
# Or run normally
go run .
```

### Code Formatting and Linting

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Vet code
go vet ./...
```

### Environment Setup

Create a `.env` file for production:

```env
PORT=8080
JWT_SECRET=your-production-secret-key
DB_URL=/app/data/chirpy.db
LOG_LEVEL=info
```

## 🔒 Security

- JWT-based authentication
- Password hashing with bcrypt
- Input validation and sanitization
- SQL injection prevention
