FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/worker .

CMD ["./worker"]
