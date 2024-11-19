# Step 1: Build the React app
FROM node:18-alpine AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Step 2: Use a lightweight Node.js image to serve the app
FROM node:18-alpine AS serve-stage
WORKDIR /app

# Install 'serve' globally
RUN npm install -g serve

# Copy the build output from the previous stage
COPY --from=build-stage /app/build ./build

# Expose the port Serve will run on
EXPOSE 3000

# Command to start the server
CMD ["serve", "-s", "build", "-l", "3000"]
