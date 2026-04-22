FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /pack-wise .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /pack-wise /usr/local/bin/pack-wise
COPY internal/db/migrations ./internal/db/migrations

EXPOSE 8080

CMD ["pack-wise"]
