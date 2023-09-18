# Use an official nodejs runtime as a parent image
FROM node:18 as build

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json
COPY package*.json ./

# Install project dependencies
RUN npm install

# Copy the rest of your application source code
COPY . .

RUN npm run build

FROM node:18-alpine

WORKDIR /app

COPY --from=build /app/package*.json ./
COPY --from=build /app/node_modules ./node_modules/
COPY --from=build /app/dist ./dist/

# Expose the port your GraphQL server will listen on
EXPOSE 3000

# Define the command to run your GraphQL server
CMD ["node", "dist/index.js"]
