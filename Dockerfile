# Build binary
FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o beli-mang-be cmd/main.go

# Build final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/beli-mang-be /app
COPY local_configuration/config.yaml /app/local_configuration/
EXPOSE 8080
CMD ["./beli-mang-be"]