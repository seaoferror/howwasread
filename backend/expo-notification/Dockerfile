# Stage 1: Build the TypeScript code
FROM node:25-alpine AS stage1
WORKDIR /app

# 1. Install build tools for native modules
RUN apk add --no-cache python3 make g++

# Copy package files and install ALL dependencies
COPY package*.json ./
RUN npm ci

# Copy the rest of the source code and run the build script
COPY . .
RUN npm run build

FROM node:25-alpine
WORKDIR /app
ENV NODE_ENV=production

# 2. Install build tools here as well, because npm ci --omit=dev
# still has to compile the Kafka production dependency
RUN apk add --no-cache python3 make g++

# Copy package files and install ONLY production dependencies
COPY package*.json ./
RUN npm ci --omit=dev

# Copy only the compiled JavaScript from the builder stage
COPY --from=stage1 /app/dist ./dist

# Run as a secure, non-root user (provided automatically by the Node image)
USER node

CMD ["node", "dist/index.js"]