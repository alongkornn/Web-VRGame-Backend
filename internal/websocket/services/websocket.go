package services

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ตัวแปลง WebSocket Connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[*websocket.Conn]bool)
var Mutex = &sync.Mutex{}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// เพิ่ม Client ใหม่
	Mutex.Lock()
	Clients[conn] = true
	Mutex.Unlock()

	// ลบ Client เมื่อปิดการเชื่อมต่อ
	defer func() {
		Mutex.Lock()
		delete(Clients, conn)
		Mutex.Unlock()
	}()
}
