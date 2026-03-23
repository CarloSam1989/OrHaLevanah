FROM golang:1.26.1 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/web /app/web
COPY --from=builder /app/data /app/data

ENV PORT=8080

EXPOSE 8080

CMD ["/app/server"]