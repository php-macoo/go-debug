package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "服务监听地址")

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()

	// WebSocket 接入点
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// 静态 HTML 客户端页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websocket/static/index.html")
	})

	log.Printf("WebSocket 聊天室已启动，访问 http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
