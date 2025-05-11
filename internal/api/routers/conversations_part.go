package routers_part

import (
	"go-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupConversationsRoutesPart(route fiber.Router) {

	route.Get("/all", service.GetConversations)
	route.Get("/user-conversation", service.GetUserConversation)
	route.Get("/messages/:conversation_id", service.GetMessagesInConversation)
	route.Get("/user-conversation", service.GetUserConversation)
	route.Get("/message-details-info/:message_id", service.GetMessageDetails)

}
