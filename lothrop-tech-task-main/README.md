# Lothrop Tech Task

A full-stack company management application with a Go backend, React TypeScript frontend, and PostgreSQL database. The system provides complete CRUD operations for company records with a modern web interface.

## Architecture

```
lothrop-tech-task/
├── backend/                # Go backend service
│   ├── cmd/server/        # Main application entry point
│   ├── internal/          # Internal packages (handlers, database, models)
│   ├── api/               # Generated API types from OpenAPI
│   ├── migrations/        # Database migrations (Sqitch)
│   ├── openapi.yaml       # API specification
│   └── Dockerfile         # Backend development container
├── frontend/              # React TypeScript frontend
│   ├── src/               # Frontend source code with shadcn/ui components
│   ├── package.json       # Node.js dependencies
│   └── Dockerfile         # Frontend development container
├── docker-compose.yml     # Multi-service orchestration
├── Makefile              # Development commands
└── README.md             # This file
```

## Design Decisions

### API Design
- **OpenAPI Specification**: Central source of truth for API contracts, ensuring consistency between backend implementation and frontend consumption
- **Versioned APIs**: `/api/v1/` prefix allows for backward compatibility and future API evolution
- **Standard HTTP Methods**: RESTful design with proper status codes and response patterns
- **UUID Primary Keys**: Better for distributed systems and reduces enumeration attacks

### Backend Architecture
- **Chi Router**: Lightweight, composable HTTP router with middleware support
- **Structured Logging**: Zap logger provides JSON output for easy monitoring and filtering
- **Code Generation**: oapi-codegen generates type-safe Go structs from OpenAPI spec
- **Database Migrations**: Sqitch provides version-controlled, declarative schema management
- **Containerization**: Docker ensures consistent development and deployment environments

### Frontend Architecture
- **React 18 + TypeScript**: Type safety and modern React features (concurrent features, Suspense)
- **shadcn/ui Components**: Consistent, accessible UI components built on Radix primitives
- **Form Validation**: React Hook Form with Zod for type-safe client-side validation
- **Hot Module Replacement**: Vite provides fast development feedback loop

### Database Design
- **PostgreSQL**: ACID compliance, robust JSON support, and excellent performance
- **Constraint Validation**: Database-level constraints ensure data integrity
- **Indexed Columns**: Strategic indexing on search and filter columns
- **Audit Fields**: Created/updated timestamps for record tracking

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Make (optional, for convenience commands)
- Sqitch (for database migrations)

### 1. Start All Services
```bash
# Clone the repository and navigate to it
cd lothrop-tech-task

# Start all services
docker compose up -d

# Or use the Makefile
make up
```

### 2. Run Database Migrations
```bash
# Navigate to migrations directory
cd backend/migrations

# Deploy migrations
sqitch deploy --verify

# Or use the Makefile
make migrate
```

### 3. Access Services

- **Frontend**: http://localhost:5174 (Company management interface)
- **Backend API**: http://localhost:8080/api/v1/ (RESTful API)
- **Database**: localhost:5432 (PostgreSQL with sample data)

## Services

### Backend Service (Port 8080)
- **Language**: Go 1.23
- **Framework**: Chi router with middleware
- **Features**:
  - Complete CRUD operations for companies
  - OpenAPI 3.0 specification with auto-generated types
  - Structured JSON logging with Zap
  - Hot reloading in development
  - Request ID tracking and CORS support

**API Endpoints:**
- `GET /api/v1/companies` - List companies with pagination
- `POST /api/v1/companies` - Create new company
- `GET /api/v1/companies/{id}` - Get company by ID
- `PUT /api/v1/companies/{id}` - Update company
- `DELETE /api/v1/companies/{id}` - Delete company
- `GET /health` - Health check endpoint

### Frontend Service (Port 5174)
- **Framework**: React 18 with TypeScript
- **UI Library**: shadcn/ui components with Tailwind CSS
- **Build Tool**: Vite with hot module replacement
- **Features**:
  - Company listing with responsive table
  - Create/edit forms with validation
  - Delete confirmation dialogs
  - Modern, accessible UI components

### Database Service (Port 5432)
- **Database**: PostgreSQL 15
- **Migration Tool**: Sqitch for version control
- **Features**:
  - Declarative schema migrations
  - UUID primary keys with proper constraints
  - Indexed columns for performance
  - Sample data for development

## Database Schema

### Companies Table
```sql
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jurisdiction VARCHAR(20) NOT NULL CHECK (jurisdiction IN ('UK', 'Singapore', 'Caymens')),
    company_name VARCHAR(255) NOT NULL,
    company_address TEXT NOT NULL,
    nature_of_business TEXT,
    number_of_directors INTEGER CHECK (number_of_directors >= 0),
    number_of_shareholders INTEGER CHECK (number_of_shareholders >= 0),
    sec_code VARCHAR(50),
    date_created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- Primary key on `id`
- Index on `jurisdiction` for filtering
- Index on `company_name` for searching
- Index on `date_created` for sorting

## Development

### Available Make Commands
```bash
make help              # Show available commands
make up                # Start all services
make down              # Stop all services
make build             # Build all services
make logs              # Show logs for all services
make logs-backend      # Show backend logs
make logs-frontend     # Show frontend logs
make logs-postgres     # Show postgres logs
make migrate           # Run database migrations
make migrate-status    # Check migration status
make test-backend      # Test backend API
make clean             # Clean up containers and volumes
make dev               # Start development environment
```

### Development Workflow

1. **Backend Development**:
   - Edit Go files in `backend/`
   - Air automatically rebuilds and restarts the server
   - Use structured logging for debugging
   - Update OpenAPI spec for new endpoints

2. **Frontend Development**:
   - Edit React/TypeScript files in `frontend/src/`
   - Vite provides instant hot module replacement
   - Use shadcn/ui components for consistency
   - Forms automatically validate with Zod schemas

3. **Database Changes**:
   - Create new migration: `sqitch add migration_name`
   - Edit migration files in `backend/migrations/`
   - Deploy: `sqitch deploy`
   - Verify changes and update API models

4. **API Changes**:
   - Update `backend/openapi.yaml` specification
   - Regenerate Go types: `go generate ./...`
   - Test endpoints with curl or API client

## Configuration

### Environment Variables

**Backend:**
- `PORT`: Server port (default: 8080)
- `POSTGRES_HOST`: Database host (default: postgres)
- `POSTGRES_PORT`: Database port (default: 5432)
- `POSTGRES_DB`: Database name (default: lothrop_db)
- `POSTGRES_USER`: Database user (default: postgres)
- `POSTGRES_PASSWORD`: Database password (default: password)

**Frontend:**
- `VITE_API_URL`: Backend API URL (default: http://localhost:8080)

**Database:**
- `POSTGRES_DB`: Database name
- `POSTGRES_USER`: Database user
- `POSTGRES_PASSWORD`: Database password

## Testing

### Backend API Testing
```bash
# Test company listing
curl http://localhost:8080/api/v1/companies

# Create a new company
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "jurisdiction": "UK",
    "company_name": "Test Company Ltd",
    "company_address": "123 Test Street, London, UK",
    "nature_of_business": "Software Development",
    "number_of_directors": 2,
    "number_of_shareholders": 3,
    "sec_code": "TEST123"
  }'

# Test health endpoint
curl http://localhost:8080/health
```

### Database Testing
```bash
# Connect to database
docker compose exec postgres psql -U postgres -d lothrop_db

# List tables
\dt

# Describe companies table
\d companies

# View sample data
SELECT * FROM companies LIMIT 5;
```

## Troubleshooting

### Common Issues

1. **Port already in use**:
   ```bash
   # Check what's using the port
   sudo lsof -i :8080
   sudo lsof -i :5173
   sudo lsof -i :5432
   
   # Stop existing containers
   docker compose down
   ```

2. **Database connection issues**:
   ```bash
   # Check if PostgreSQL is healthy
   docker compose ps
   
   # View PostgreSQL logs
   docker compose logs postgres
   ```

3. **Frontend build issues**:
   ```bash
   # Rebuild frontend container
   docker compose build frontend --no-cache
   docker compose up -d frontend
   ```

4. **Migration issues**:
   ```bash
   # Check migration status
   cd backend/migrations && sqitch status
   
   # Revert last migration
   sqitch revert --to @HEAD^
   ```

## Features

This application provides a complete company management system with:

- **Company CRUD Operations**: Create, read, update, and delete company records
- **Data Validation**: Client-side and server-side validation with proper error handling
- **Responsive Design**: Modern UI that works on desktop and mobile devices
- **Type Safety**: Full TypeScript coverage from database to frontend
- **Development Environment**: Hot reloading and containerized services
- **Database Migrations**: Version-controlled schema with sample data
- **API Documentation**: OpenAPI specification for clear API contracts
- **Structured Logging**: JSON logs for monitoring and debugging
- **Error Handling**: Proper HTTP status codes and error messages

## Technical Stack

- **Backend**: Go 1.23, Chi router, PostgreSQL, Zap logging
- **Frontend**: React 18, TypeScript, Vite, shadcn/ui, Tailwind CSS
- **Database**: PostgreSQL 15 with Sqitch migrations
- **Development**: Docker Compose, Air (Go hot reload), Vite HMR
- **API**: OpenAPI 3.0 specification with code generation

## License

This project is for development purposes as part of the Lothrop Tech Task.
