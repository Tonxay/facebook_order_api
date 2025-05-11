package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	gormpkg "github.com/yourusername/go-api/internal/pkg"
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
		db := gormpkg.GetDB()
		// Optional: Set a Pong handler to keep connection alive
		c.SetPongHandler(func(appData string) error {
			fmt.Println("Pong received")
			return nil
		})

		// Read loop (optional but helps keep connection alive)
		go func() {
			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					fmt.Println("WebSocket read error:", err)
					c.Close()
					break
				}
			}
		}()

		// Write loop
		for {
			data, err := dbservice.Getcustomers(db)
			if err != nil {
				fmt.Println("DB error:", err)
				return
			}

			result, err := json.Marshal(data)
			if err != nil {
				fmt.Println("JSON marshal error:", err)
				return
			}

			err = c.WriteMessage(websocket.TextMessage, result)
			if err != nil {
				fmt.Println("WebSocket write error:", err)
				return
			}

			// Wait 3 seconds before sending again to avoid spamming
			time.Sleep(3 * time.Second)
		}
	}
}
