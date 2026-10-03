# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/gateway cmd/server/main.go

# Production stage
FROM alpine:3.20

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/gateway /app/gateway

EXPOSE 8080
ENTRYPOINT ["/app/gateway"]
