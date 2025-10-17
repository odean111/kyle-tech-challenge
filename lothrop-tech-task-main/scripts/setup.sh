#!/bin/bash

# Setup script for the Lothrop Tech Task
echo "Setting up Lothrop Tech Task..."

# Start the services
echo "Starting services with docker-compose..."
docker-compose up -d postgres

# Wait for postgres to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 10

# Run migrations
echo "Running database migrations..."
cd backend/migrations
docker-compose exec postgres sh -c "PGPASSWORD=password psql -h localhost -U postgres -d lothrop_db -c '\l'"

echo "Setup complete! You can now run 'docker-compose up' to start all services."