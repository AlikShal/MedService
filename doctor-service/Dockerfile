FROM golang:1.21-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/doctor-service ./cmd/doctor-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /bin/doctor-service /app/doctor-service

EXPOSE 8080

ENTRYPOINT ["/app/doctor-service"]
