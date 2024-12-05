# Stage 1: Build Tailwind CSS
FROM node:16 AS tailwind-build

WORKDIR /app

# Copy package.json and install npm dependencies
COPY package*.json ./
RUN npm install

# Copy the Tailwind config and source CSS file
COPY tailwind.config.js ./
COPY input.css ./
COPY static ./static
COPY templates ./templates

# Build the Tailwind CSS file
RUN npm run tailwind

# Output the generated CSS file to a known location
RUN mkdir -p /output/static/css && cp -rf ./static/css /output/static/

# Inject CSS file hash into the HTML template to prevent caching across changes
RUN HASH=$(sha256sum /output/static/css/styles.css | awk '{print $1}') && \
    sed -i "s|{HASH_PLACEHOLDER}|${HASH}|" ./templates/base.html && \
    cp -rf ./templates /output/templates

RUN ls -la /output/static


# Stage 2: Build Go application
FROM golang:1.23 AS go-build

WORKDIR /app

# Copy Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go source code
COPY . .

# Build the Go binary (statically linked for deployment)
RUN CGO_ENABLED=1 GOOS=linux go build -tags "sqlite_fts5" -o /output/app .


# Stage 3: Final runtime container
FROM debian:bookworm-slim

WORKDIR /app

# Copy the generated Tailwind CSS, modified themplates, and Go binary from the previous stages
COPY --from=tailwind-build /output/static ./static
COPY --from=tailwind-build /output/templates ./templates
COPY --from=go-build /output/app ./app
COPY migrations ./migrations

# Expose the port that the Go app will listen on (default 8080)
EXPOSE 8080

# Run the Go binary when the container starts
CMD ["./app"]
