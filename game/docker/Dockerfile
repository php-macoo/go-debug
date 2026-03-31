# 构建上下文为仓库根目录：docker build -f game/docker/Dockerfile .
FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/game-svc ./game/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/game-svc /app/game-svc
COPY game/static /app/game/static
EXPOSE 8082
USER nobody:nobody
ENTRYPOINT ["/app/game-svc"]
CMD ["-config", "/config/config.yaml"]
