FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /users-api ./cmd/server

FROM alpine:3.20

RUN adduser -D appuser
WORKDIR /app

COPY --from=builder /users-api /users-api

USER appuser
EXPOSE 8080

CMD ["/users-api"]
