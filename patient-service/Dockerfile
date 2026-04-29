FROM golang:1.21-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/patient-service ./cmd/patient-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /bin/patient-service /app/patient-service

EXPOSE 8083

ENTRYPOINT ["/app/patient-service"]
