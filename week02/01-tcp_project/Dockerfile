# 使用 Go 作为构建阶段的基础镜像
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod . 
COPY go.sum .
RUN go mod download

# 复制项目的源代码
COPY . .

# 定义构建目标
ARG TARGET

# 根据目标构建相应的二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tcp_app ./cmd/tcp_${TARGET}.go

# 使用 Ubuntu 20.04 作为运行时基础镜像
FROM ubuntu:20.04

# 创建无家目录的用户并设置权限，创建/opt目录
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    addgroup --system appgroup && \
    adduser --system --ingroup appgroup --no-create-home appuser && \
    mkdir -p /opt && \
    chown appuser:appgroup /opt

# 设置工作目录为/opt
WORKDIR /opt

# 复制构建的二进制文件
ARG TARGET
COPY --from=builder /tcp_app /opt/tcp_app

# 更改文件权限以使非 root 用户可以执行它
RUN chown -R appuser:appgroup /opt && \
    chmod +x /opt/tcp_app

# 切换到非 root 用户
USER appuser

# 启动相应的应用
CMD ["/opt/tcp_app"]

