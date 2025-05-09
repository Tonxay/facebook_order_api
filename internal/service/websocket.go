package service

import (
	"encoding/json"
	"sync"

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
