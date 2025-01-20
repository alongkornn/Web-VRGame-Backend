package websocket_services

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// ตัวแปลง WebSocket Connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Clients map ที่เก็บ WebSocket connections
var Clients = make(map[*websocket.Conn]bool)
var Mutex = &sync.Mutex{}

// HandleWebSocket ใช้กับ Echo
func HandleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return err
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

	// รอรับข้อความจาก WebSocket
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			return err
		}
	}
}
