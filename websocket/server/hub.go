package main

import "log"

// Hub 管理所有活跃的客户端连接，负责广播消息
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[Hub] 客户端 %s 已连接，当前在线: %d", client.id, len(h.clients))
			h.broadcastSystemMsg([]byte("系统: " + client.id + " 加入了聊天室"))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[Hub] 客户端 %s 已断开，当前在线: %d", client.id, len(h.clients))
				h.broadcastSystemMsg([]byte("系统: " + client.id + " 离开了聊天室"))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 发送缓冲区满，移除该客户端
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) broadcastSystemMsg(msg []byte) {
	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
