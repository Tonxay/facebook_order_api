package service

import (
	"encoding/json"
	"fmt"
	gormpkg "go-api/internal/pkg"
	dbservice "go-api/internal/service/db_service"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
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

// /
var ClientALLs = make(map[*websocket.Conn]bool)
var LockALLs sync.Mutex

func RegisterClientAll(conn *websocket.Conn) {
	LockALLs.Lock()
	ClientALLs[conn] = true
	LockALLs.Unlock()
}

func RemoveClientAll(conn *websocket.Conn) {
	LockALLs.Lock()
	delete(ClientALLs, conn)
	LockALLs.Unlock()
}

func PushToAll(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	LockALLs.Lock()
	defer LockALLs.Unlock()

	for conn := range ClientALLs {
		if conn == nil {
			delete(ClientALLs, conn)
			continue
		}

		// Set write deadline to avoid hanging
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error sending to client, removing connection:", err)
			conn.Close()
			delete(ClientALLs, conn)
		}
	}
}

// func PushToAll(payload interface{}) {
// 	data, err := json.Marshal(payload)

// 	if err != nil {
// 		log.Println("JSON marshal error:", err)
// 		return
// 	}
// 	LockALLs.Lock()
// 	defer LockALLs.Unlock()

//		for conn := range ClientALLs {
//			err := conn.WriteMessage(websocket.TextMessage, data)
//			if err != nil {
//				log.Println("Error sending to client:", err)
//				conn.Close()
//				delete(ClientALLs, conn)
//			}
//		}
//	}
// func WebSocketMessageAllUserHandler() func(*websocket.Conn) {
// 	return func(c *websocket.Conn) {
// 		// Register connection
// 		RegisterClientAll(c)
// 		defer func() {
// 			RemoveClientAll(c)
// 			c.Close()
// 		}()

// 		// Set pong handler
// 		c.SetPongHandler(func(appData string) error {
// 			log.Println("Pong received")
// 			return nil
// 		})

// 		// Heartbeat: send pings every 30 seconds
// 		go func(conn *websocket.Conn) {
// 			ticker := time.NewTicker(30 * time.Second)
// 			defer ticker.Stop()

// 			for {
// 				<-ticker.C
// 				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
// 					log.Println("Ping failed:", err)
// 					conn.Close()
// 					RemoveClientAll(conn)
// 					break
// 				}
// 			}
// 		}(c)

// 		// Read loop (to keep connection alive)
// 		for {
// 			if _, _, err := c.ReadMessage(); err != nil {
// 				log.Println("Read error:", err)
// 				break
// 			}
// 		}
// 	}
// }

func WebSocketMessageAllUserHandler() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// userID := c.Query("user_id")
		// if userID == "" {
		// 	log.Println("Missing user_id query")
		// 	c.Close()
		// 	return
		// }

		RegisterClientAll(c)

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}

		RemoveClientAll(c)

	}
}

func WebSocketMessageHandler() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		userID := c.Query("user_id")
		if userID == "" {
			log.Println("Missing user_id query")
			c.Close()
			return
		}
		RegisterClient(userID, c)
		log.Printf("Connected WebSocket for user: %s", userID)

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}

		RemoveClient(userID)
		log.Printf("Disconnected WebSocket for user: %s", userID)
	}
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
