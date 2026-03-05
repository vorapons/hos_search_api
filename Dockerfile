# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/

# ── Stage 2: Run ──────────────────────────────────────────────────────────────
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3458

CMD ["./server"]
