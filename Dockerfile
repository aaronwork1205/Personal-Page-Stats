
# ---- Stage 1: Build ----
FROM golang:1.22-alpine AS builder

# 安装 SQLite 依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# 复制代码并构建
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

# ---- Stage 2: Run ----
FROM alpine:latest

RUN apk add --no-cache sqlite

WORKDIR /app

# 拷贝二进制和 db 文件夹
COPY --from=builder /app/server .
COPY --from=builder /app/data ./data

EXPOSE 8080

CMD ["./server"]
