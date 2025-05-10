# Build executable binary
FROM golang:1.23.2-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN apk add --no-cache ca-certificates git tzdata && go mod tidy

COPY . .

RUN go build -ldflags="-s -w -extldflags '-static'" -installsuffix cgo -o /bin/api-app ./cmd/main/main.go

# Use alpine image as runtime
FROM alpine:3.16 AS release

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/api-app /bin/api-app

# Runtime environment variables (can be overwritten when running `docker run`)
ENV API_VERSION=dev \
    BUILD_DATE=unknown

# Command to run the binary
ENTRYPOINT ["/bin/api-app"]
