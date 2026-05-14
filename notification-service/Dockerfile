FROM golang:1.21-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notification-service ./cmd/notification-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /bin/notification-service /app/notification-service

EXPOSE 8085

ENTRYPOINT ["/app/notification-service"]
