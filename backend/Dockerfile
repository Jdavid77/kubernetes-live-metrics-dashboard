# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
# Enable automatic Go toolchain download for newer versions
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOTOOLCHAIN=auto go build -o api ./cmd/api

# Stage 2: Runtime
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

COPY --from=builder /app/api .

USER 65534:65534

EXPOSE 8080

CMD ["./api"]
