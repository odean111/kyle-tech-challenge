FROM node:20-alpine AS development

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 5173

# Command to run vite dev server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]