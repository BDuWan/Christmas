# Stage 1: Build binary
FROM golang:1.21-bookworm AS builder

WORKDIR /app

# Copy go.mod (và nếu có thì cả go.sum) để cache dependency
COPY go.mod ./
# Nếu sau này bạn có go.sum, bỏ comment dòng dưới:
# COPY go.sum ./

# Tải dependency (sẽ tạo go.sum bên trong image)
RUN go mod download

# Copy toàn bộ source (Go + HTML + img + music)
COPY . .

# Đảm bảo go.sum đầy đủ trước khi build
RUN go mod tidy

# Build binary
RUN go build -o server server.go

# Stage 2: Runtime (Ubuntu)
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary + static files
COPY --from=builder /app/server /app/server
COPY --from=builder /app/1.html /app/1.html
COPY --from=builder /app/2.html /app/2.html
COPY --from=builder /app/3.html /app/3.html
COPY --from=builder /app/img /app/img
COPY --from=builder /app/music /app/music

EXPOSE 8080

CMD ["./server"]

