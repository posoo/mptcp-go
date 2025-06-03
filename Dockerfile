FROM golang:1.24-alpine3.22 AS builder

WORKDIR /app

COPY . .
RUN go build -o server ./cmd/webserver

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8684

CMD ["./server", "8684"]