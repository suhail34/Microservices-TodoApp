# Use an official nodejs runtime as a parent image
FROM node:18 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json
COPY package*.json ./

# Install project dependencies
RUN npm install

# Copy the rest of your application source code
COPY . .

# Expose the port your GraphQL server will listen on
EXPOSE 8080

# Define the command to run your GraphQL server
CMD ["npm", "start"]
