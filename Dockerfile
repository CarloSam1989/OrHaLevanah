# Etapa de build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Etapa final
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/web /app/web
COPY --from=builder /app/data /app/data

ENV PORT=8080

EXPOSE 8080

CMD ["/app/server"]