package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/websocket/v2"
	dbservice "github.com/yourusername/go-api/internal/service/db_service"
)

var Clients = make(map[string]*websocket.Conn)
var Lock sync.Mutex

func RegisterClient(userID string, conn *websocket.Conn) {
	Lock.Lock()
	Clients[userID] = conn
	Lock.Unlock()
}

func RemoveClient(userID string) {
	Lock.Lock()
	delete(Clients, userID)
	Lock.Unlock()
}

func PushToUser(userID string, payload interface{}) {
	Lock.Lock()
	conn, ok := Clients[userID]
	Lock.Unlock()
	if !ok {
		return
	}

	data, _ := json.Marshal(payload)
	conn.WriteMessage(websocket.TextMessage, data)
}

func PurchaseWebSocketCheckPayment() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()

		// db := gormpkg.GetDB()

		for {
			// _, msg, err := c.ReadMessage()
			// if err != nil {
			// 	fmt.Println("WebSocket Read error:", err)
			// 	return
			// }

			// Validate UUID
			// if err := validators.ValidateUuId(qrPaymentId); err != nil {
			// 	writeAndClose(c, "Invalid UUID")
			// 	return
			// }

			data, err := dbservice.Getcustomers()
			if err != nil {
				// writeAndClose(c, data)
				return
			}
			resulf, _ := json.Marshal(data)
			// Successful response
			if err := c.WriteMessage(websocket.TextMessage, []byte(resulf)); err != nil {
				fmt.Println("WebSocket Write error:", err)
				return
			}
		}
	}
}
