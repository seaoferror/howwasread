# Stage 1: Build the TypeScript code
FROM node:24-alpine AS stage1
WORKDIR /app

# Copy package files and install ALL dependencies (including dev tools like typescript, tsc-alias)
COPY package*.json ./
RUN npm ci

# Copy the rest of the source code and run the build script
COPY . .
RUN npm run build

FROM node:24-alpine
WORKDIR /app
ENV NODE_ENV=production

# Copy package files and install ONLY production dependencies
COPY package*.json ./
RUN npm ci --omit=dev

# Copy only the compiled JavaScript from the builder stage
COPY --from=stage1 /app/dist ./dist

# Run as a secure, non-root user (provided automatically by the Node image)
USER node

CMD ["node", "dist/index.js"]