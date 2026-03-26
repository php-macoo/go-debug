# go-debug

Go 工具、脚本、协议 Demo 的集合仓库，用于学习与调试各类技术点。

## 项目结构

```
go-debug/
├── websocket/          # WebSocket Demo：广播聊天室
│   ├── server/
│   │   ├── main.go     # 入口：HTTP 服务 + 路由
│   │   ├── hub.go      # Hub：连接管理 & 广播
│   │   └── client.go   # Client：读写泵 + ping/pong
│   └── static/
│       └── index.html  # 浏览器聊天室页面
├── http/               # (待扩展) HTTP 工具 Demo
├── grpc/               # (待扩展) gRPC Demo
└── tools/              # (待扩展) 各类工具脚本
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
