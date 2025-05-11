package routers_part

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/yourusername/go-api/internal/service"
)

func SetupWebSocketRoutesPart(app *fiber.App) {
	// WebSocket route
	// app.Get("/ws/:facebook_id", websocket.New(func(c *websocket.Conn) {
	// 	userID := c.Params("facebook_id")
	// 	service.RegisterClient(userID, c)

	// 	defer func() {
	// 		service.RemoveClient(userID)
	// 		c.Close()
	// 	}()

	// 	for {
	// 		if _, _, err := c.ReadMessage(); err != nil {
	// 			break
	// 		}
	// 	}
	// }))
	app.Get("/ws/customers", websocket.New(service.PurchaseWebSocketCheckPayment()))

}
