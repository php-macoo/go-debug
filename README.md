# go-debug

Go 工具、脚本、协议 Demo 的集合仓库，用于学习与调试各类技术点。

## 项目结构

```
go-debug/
├── websocket/              # WebSocket Demo：广播聊天室
│   ├── server/
│   │   ├── main.go
│   │   ├── hub.go
│   │   └── client.go
│   └── static/
│       └── index.html
├── game/                   # 小游戏空间：Gin + MySQL + 静态前端
│   ├── main.go             # 入口
│   ├── config.yaml         # 服务 / 数据库 / 认证 / 成绩规则（勿提交真实密钥）
│   ├── config/             # 配置加载
│   ├── model/              # GORM 模型
│   ├── dao/                # 数据访问
│   ├── service/            # 业务逻辑
│   ├── handler/            # HTTP 处理
│   ├── middleware/         # 认证、访问日志
│   ├── router/             # 路由
│   ├── pkg/                # resp、工具
│   ├── docker/             # Dockerfile、compose、config.docker.yaml
│   ├── script/             # 运维脚本（SQL、docker-deploy.sh）
│   │   └── sql/schema.sql  # 基准建表（人工执行，应用内不跑 AutoMigrate）
│   └── static/             # 大厅与各游戏页面
│       ├── index.html
│       ├── css/
│       ├── js/common/      # auth.js（含匿名 X-Guest-Device-Id）
│       └── games/          # match3、2048、snake 等
├── http/                   # (待扩展) HTTP 工具 Demo
├── grpc/                   # (待扩展) gRPC Demo
└── tools/                  # (待扩展) 各类工具脚本
```

## WebSocket Demo

### 功能

- 多客户端同时连接，消息实时广播
- 系统消息：用户加入 / 离开通知
- Ping/Pong 心跳保活（60s 超时检测）
- 批量刷写：积压消息合并发送，减少系统调用
- 带颜色区分的浏览器聊天室 UI（自己/他人/系统）

### 运行

```bash
# 在项目根目录执行
go run ./websocket/server/

# 浏览器打开（可多标签模拟多用户）
open http://localhost:8080
```

### 架构

```
Browser ──ws──► Client.readPump  ──► Hub.broadcast ──► 所有 Client.writePump ──ws──► Browser
                Client.writePump ◄──
                     ▲ Ping/Pong ticker
```

| 组件 | 职责 |
|------|------|
| `Hub` | 管理连接集合；处理注册/注销；广播消息 |
| `Client` | 封装单条 WebSocket 连接；独立 goroutine 读写 |
| `readPump` | 从浏览器读消息 → 投递 Hub |
| `writePump` | 从 Hub 收消息 → 写回浏览器；发 Ping 心跳 |

## 小游戏空间（game）

基于 **Gin + GORM + MySQL** 的网页小游戏大厅：用户注册/登录、多游戏入口、成绩上报与排行榜、API 访问日志。静态资源由 Gin `NoRoute` 回退到本地目录。

### 功能概览

- **大厅**：`/api/games` 列出已上线游戏；首页 `static/index.html` 跳转各游戏页。
- **认证**：`/api/auth/register`、`/api/auth/login`；`/api/user/*` 需 `Authorization: Bearer <token>`。
- **游玩与成绩**（`/api/game/:gameKey/...`）：
  - 可选登录：有 token 用登录用户；无 token 时前端自动带 **`X-Guest-Device-Id`**（`localStorage` 持久化），服务端为匿名用户创建 `users` 中 `source=guest` 记录，逻辑与正式用户一致。
  - `POST .../run/start` 领取对局凭证 `runId`；`POST .../score` 提交成绩（服务端校验游戏上线、完成时间上下限、`runId`、提交间隔等）。
  - `GET .../leaderboard` 公开，无需登录。
- **表结构**：以 `game/script/sql/schema.sql` 为准在库中执行；**不**在程序内 AutoMigrate，后续改表由维护者提供增量 SQL 手工执行。
- **开发调试**：`config.yaml` 中 `database.log_sql: true` 可打开 GORM SQL 日志（生产请关闭）。

### 环境要求

- Go 1.20+
- MySQL 5.7+ / 8.x

### 初始化数据库

```bash
# 在目标库中执行（可先 CREATE DATABASE）
mysql -h127.0.0.1 -P3308 -uroot -p game_db < game/script/sql/schema.sql
```

按本机修改 `game/config.yaml` 里的 `database` 段（host、port、user、password、name）。

### 指定配置文件路径

默认读取当前工作目录下的 **`game/config.yaml`**。可用启动参数 **`-config`** 指向其它文件（适合挂载到固定路径的场景）。

```bash
go run ./game/ -config /etc/game/config.yaml

# 或相对路径
go run ./game/ -config ./my-config.yaml
```

`server.static_dir` 等与路径相关的项仍相对**进程当前工作目录**，除非写成绝对路径。

### Docker 部署

- **`game/docker/Dockerfile`**：多阶段构建，镜像内包含二进制与 `game/static`，默认 **`CMD` 为 `-config /config/config.yaml`**（需挂载配置文件，镜像内可不打包 yaml）。
- **`game/docker/docker-compose.yml`**：仅启动 **`app`**，**不包含数据库**。请自备 MySQL，并先在目标库执行 **`game/script/sql/schema.sql`**。编辑 **`game/docker/config.docker.yaml`** 中的 `database`（示例使用 `host.docker.internal` 访问宿主机上的 MySQL；云库改为域名/IP）。Compose 中为 Linux 增加了 `extra_hosts: host.docker.internal:host-gateway`。
- **`game/script/docker-deploy.sh`**：在仓库根执行 compose（`--project-directory` 指向根目录）。

```bash
cd /path/to/go-debug
bash game/script/docker-deploy.sh up    # 构建并后台启动
bash game/script/docker-deploy.sh logs  # 看 app 日志
bash game/script/docker-deploy.sh down  # 停止容器
```

访问 `http://localhost:8082`。生产环境务必修改 **`config.docker.yaml`** 中的 `auth.token_secret` 与数据库账号密码。

单独构建镜像（上下文为仓库根）：

```bash
docker build -f game/docker/Dockerfile -t game-space:local .
docker run --rm -p 8082:8082 -v "$PWD/game/docker/config.docker.yaml:/config/config.yaml:ro" game-space:local
```

### 配置说明

| 配置块 | 说明 |
|--------|------|
| `server` | `addr` 监听地址；`static_dir` 相对**项目根**的静态目录（默认 `game/static`）。 |
| `database` | MySQL 连接；`log_sql` 是否打印 SQL。 |
| `auth` | `token_secret`、`token_expire_days`。 |
| `score` | 提交最小间隔、`run` 凭证 TTL、默认与各 `gameKey` 的完成时间毫秒上下限。 |

### 运行

```bash
# 仓库根目录
go run ./game/

# 浏览器（端口以 config.yaml 为准，默认 8082）
open http://localhost:8082
```

### 三消（match3）前端特性

- 类「羊了个羊」叠层消除：暂存槽凑 3 张相同消除，多层遮挡不可点。
- 道具：撤回、洗牌、移出；Web Audio 音效；胜利彩带与计时。
- 未登录也可完整游玩并上榜（匿名设备 ID）；登录后使用独立账号成绩。

### API 约定（自定义客户端）

- JSON 字段与前端一致：**小驼峰**（如 `completionTimeMs`）。
- 匿名游玩：`run/start` 与 `score` 须在无 Bearer 时携带请求头 **`X-Guest-Device-Id`**（建议 UUID，长度 8～128）。
