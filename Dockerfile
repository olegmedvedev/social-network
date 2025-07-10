# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/gateway

FROM alpine:3.22
WORKDIR /root/
COPY --from=builder /app/server .
COPY .env .
EXPOSE 8080
CMD ["./server"] 