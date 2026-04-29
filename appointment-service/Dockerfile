FROM golang:1.21-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/appointment-service ./cmd/appointment-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /bin/appointment-service /app/appointment-service

EXPOSE 8081

ENTRYPOINT ["/app/appointment-service"]
