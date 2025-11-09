# # Build executable binary
# FROM golang:1.23.2-alpine AS builder

# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# WORKDIR /build

# COPY go.mod .
# COPY go.sum .
# RUN apk add --no-cache ca-certificates git tzdata && go mod tidy

# COPY . .

# RUN go build -ldflags="-s -w -extldflags '-static'" -installsuffix cgo -o /bin/api-app ./cmd/main/main.go

# # Use alpine image as runtime
# FROM alpine:3.16 AS release


# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /bin/api-app /bin/api-app

# # Copy font file (adjust path as needed)
# COPY Phetsarath_OT.ttf /app/Phetsarath_OT.ttf

# # Runtime environment variables (can be overwritten when running `docker run`)

# WORKDIR /app
# ARG API_VERSION
# ARG BUILD_DATE
# ENV API_VERSION ${API_VERSION}
# ENV BUILD_DATE ${BUILD_DATE}



# # Command to run the binary
# ENTRYPOINT ["/bin/api-app"]

# # Build executable binary
# FROM golang:1.24.0-alpine AS builder

# # Disable CGO for static binary
# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# WORKDIR /build

# # Copy go.mod and go.sum first for caching
# COPY go.mod go.sum ./

# # Install system dependencies separately for better error visibility
# RUN apk add --no-cache ca-certificates git tzdata

# # Download dependencies
# RUN go mod tidy

# # Copy the rest of the app
# COPY . .

# RUN apt-get update && apt-get install -y chromium
# # Build the binary
# RUN go build -ldflags="-s -w -extldflags '-static'" -installsuffix cgo -o /bin/api-app ./cmd/main/main.go

# # Use alpine image as runtime

# FROM alpine:3.22 AS release

# # Copy timezone & certs from builder
# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /bin/api-app /bin/api-app

# # Copy font file
# COPY Phetsarath_OT.ttf /app/Phetsarath_OT.ttf

# WORKDIR /app

# # Correct ENV syntax
# ARG API_VERSION
# ARG BUILD_DATE
# ENV API_VERSION=${API_VERSION}
# ENV BUILD_DATE=${BUILD_DATE}

# # Run the binary
# ENTRYPOINT ["/bin/api-app"]
# # Build executable binary
FROM golang:1.24.0-alpine AS builder

# Disable CGO for static binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./

# Install system dependencies for Alpine
RUN apk add --no-cache ca-certificates git tzdata

# Download dependencies
RUN go mod tidy

# Copy the rest of the app
COPY . .

# Build the binary
RUN go build -ldflags="-s -w -extldflags '-static'" -installsuffix cgo -o /bin/api-app ./cmd/main/main.go

RUN apt-get update && apt-get install -y ca-certificates
# -----------------------------------------------------------------
# Use debian-slim as runtime, NOT alpine
# -----------------------------------------------------------------
FROM debian:bookworm-slim AS release

# Install Chromium and its dependencies
# This is the correct place to install it

# RUN apt-get update && apt-get install -y --no-install-recommends \
#     chromium \
#     ca-certificates \
#     tzdata \
#     # Add minimal dependencies for headless chromium to run
#     libasound2 \
#     libgbm1 \
#     libgtk-3-0 \
#     libxshmfence1 \
# && rm -rf /var/lib/apt/lists/*

# Copy the compiled app from the builder
COPY --from=builder /bin/api-app /bin/api-app

# Copy font file
COPY Phetsarath_OT.ttf /app/Phetsarath_OT.ttf

WORKDIR /app

# Correct ENV syntax
ARG API_VERSION
ARG BUILD_DATE
ENV API_VERSION=${API_VERSION}
ENV BUILD_DATE=${BUILD_DATE}

# Run the binary
ENTRYPOINT ["/bin/api-app"]