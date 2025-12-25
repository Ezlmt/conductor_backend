# Conductor Backend

A RESTful API backend for a course management system built with Go, Gin, and PostgreSQL.

## Features

- **User Authentication**: JWT-based authentication with role-based access control
- **Role Management**: Support for Students and Professors
- **Course Management**: Create, delete, and query courses
- **Enrollment System**: Students can join/leave courses using course codes
- **CORS Support**: Configured for frontend integration

## Tech Stack

- **Language**: Go 1.24.2
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 16
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt

## Project Structure

```
conductor_backend/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── docker-compose.yml      # PostgreSQL database setup
├── internal/
│   ├── controllers/        # Request handlers
│   │   ├── user.go        # User authentication endpoints
│   │   └── course.go      # Course management endpoints
│   ├── database/          # Database connection and configuration
│   │   └── db.go
│   ├── middleware/        # HTTP middleware
│   │   ├── auth.go        # JWT authentication middleware
│   │   ├── role.go        # Role-based access control
│   │   └── dev.go         # Development-only endpoints
│   ├── models/            # Data models
│   │   ├── user.go
│   │   ├── course.go
│   │   └── enrollment.go
│   └── routes/            # Route definitions
│       └── routes.go
└── README.md
```

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 16 (or use Docker Compose)
- Docker and Docker Compose (optional, for database setup)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd conductor_backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**

   Option A: Using Docker Compose (Recommended)
   ```bash
   docker-compose up -d
   ```

   Option B: Manual PostgreSQL setup
   - Create a PostgreSQL database named `conductor`
   - Update database credentials in `internal/database/db.go` if needed

4. **Configure environment variables**

   Create a `.env` file in the root directory:
   ```env
   JWT_SECRET=your-secret-key-here
   ```

## Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:9916`

## API Endpoints

### Public Endpoints

#### Health Check
- `GET /ping` - Health check endpoint

#### Authentication
- `POST /users/register` - Register a new user
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "role": 1  // 1 = Student, 2 = Professor
  }
  ```

- `POST /users/login` - Login and get JWT token
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

### Protected Endpoints (Require JWT Token)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

#### User
- `GET /me` - Get current user information

#### Course Management (Professor Only)
- `POST /courses` - Create a new course
  ```json
  {
    "name": "Introduction to Computer Science",
    "code": "CS101"
  }
  ```

- `DELETE /courses` - Delete a course
  ```json
  {
    "id": 1
  }
  ```

- `GET /courses` - Get all courses created by the current professor

#### Enrollment (Student Only)
- `POST /courses/join` - Join a course using course code
  ```json
  {
    "code": "CS101"
  }
  ```

- `DELETE /courses/:id/leave` - Leave a course

- `GET /enrollments` - Get all courses enrolled by the current student

### Development Endpoints

These endpoints are only available in development mode (controlled by middleware):

- `DELETE /dev/courses/:id` - Delete a course by ID
- `GET /dev/show-all-courses` - Show all courses in the system

## Database Models

### User
- `ID` (uint, primary key)
- `Email` (string, unique, not null)
- `PasswordHash` (string, not null)
- `Role` (int8, not null) - 1: Student, 2: Professor
- `CreatedAt` (time.Time)

### Course
- `ID` (uint, primary key)
- `Name` (string)
- `Code` (string, unique)
- `ProfessorID` (uint, foreign key)
- `CreatedAt` (time.Time)

### Enrollment
- `ID` (uint, primary key)
- `UserID` (uint, foreign key)
- `CourseID` (uint, foreign key)
- `CreatedAt` (time.Time)

## Configuration

### Database Configuration

Default database settings (can be modified in `internal/database/db.go`):
- Host: `localhost`
- Port: `5432`
- User: `user`
- Password: `123`
- Database: `conductor`

### CORS Configuration

CORS is configured in `main.go` to allow requests from `http://localhost:5173` (typical Vite dev server). Modify the `AllowOrigins` array to match your frontend URL.

### Server Port

The server runs on port `9916` by default. Change this in `main.go`:
```go
r.Run(":9916")
```

## Security Features

- Password hashing using bcrypt
- JWT token-based authentication
- Role-based access control (RBAC)
- CORS protection
- Input validation

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o conductor_backend
```

### Running with Docker
You can containerize the application using Docker. Create a `Dockerfile` and use Docker Compose to orchestrate both the database and the application.

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]


