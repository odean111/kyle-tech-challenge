FROM golang:1.23-alpine AS development

# Install required packages
RUN apk add --no-cache postgresql-client

# Install air for live reloading (using version compatible with Go 1.23)
RUN go install github.com/cosmtrek/air@v1.40.4

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Command to run air for development
CMD ["air", "-c", ".air.toml"]