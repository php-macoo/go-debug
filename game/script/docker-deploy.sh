#!/usr/bin/env bash
# 在仓库根目录通过 compose 仅构建并启动 app；数据库自备，需已执行 game/script/sql/schema.sql。
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE_FILE="$ROOT/game/docker/docker-compose.yml"

cd "$ROOT"

usage() {
  echo "用法: $(basename "$0") <命令>"
  echo "  up      构建镜像并后台启动（默认）"
  echo "  down    停止并删除容器"
  echo "  logs    跟踪 app 日志"
  echo "  build   仅构建 app 镜像"
  echo "  ps      查看状态"
}

cmd="${1:-up}"

case "$cmd" in
  up)
    docker compose -f "$COMPOSE_FILE" --project-directory "$ROOT" up -d --build
    echo "已启动。访问 http://localhost:8082 （数据库请在 config.docker.yaml 中配置）"
    ;;
  down)
    docker compose -f "$COMPOSE_FILE" --project-directory "$ROOT" down
    ;;
  logs)
    docker compose -f "$COMPOSE_FILE" --project-directory "$ROOT" logs -f app
    ;;
  build)
    docker compose -f "$COMPOSE_FILE" --project-directory "$ROOT" build app
    ;;
  ps)
    docker compose -f "$COMPOSE_FILE" --project-directory "$ROOT" ps
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    echo "未知命令: $cmd" >&2
    usage >&2
    exit 1
    ;;
esac
