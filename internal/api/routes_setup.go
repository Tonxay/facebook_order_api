package api

import (
	routers_part "go-api/internal/api/routers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	webhook := app.Group("/webhook")
	routers_part.SetupWebhookRoutesPart(webhook)

	conversation := app.Group("/conversations")
	routers_part.SetupConversationsRoutesPart(conversation)

	customer := app.Group("/customers")
	routers_part.SetupCustomersRoutesPart(customer)

	products := app.Group("/products")
	routers_part.SetupProductRoutesPart(products)

}

func SetupWebsocketRoutes(app *fiber.App) {
	routers_part.SetupWebSocketRoutesPart(app)
}

// curl -X POST https://pro.api.anousith.express/graphql \
//   -H "Content-Type: application/json" \
//   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
//   -H "Origin: https://app.anousith.express" \
//   -H "Referer: https://app.anousith.express/" \
//   -H "Accept-Encoding: gzip" \
//   --data '{
//     "operationName": "ItemsV2",
//     "variables": {
//       "where": {
//         "originSendDate_gte": "2025-04-01",
//         "originSendDate_lt": "2025-05-14",
//         "multipleItemStatus": [
//           "TRANSIT_TO_DEST_BRANCH",
//           "TRANSIT_TO_ORIGIN_BRANCH"
//         ],
//         "searchMultipleCOD": ["0", "1"],
//         "customerId": 6216826,
//         "isDeleted": 0
//       },
//       "orderBy": "originSendDate_DESC",
//       "skip": 0,
//       "limit": 100
//     },
//     "query": "query ItemsV2($where: ItemV2WhereInput, $skip: Int, $noLimit: Boolean, $limit: Int, $orderBy: OrderByItem) {\n  itemsV2(\n    where: $where\n    skip: $skip\n    noLimit: $noLimit\n    limit: $limit\n    orderBy: $orderBy\n  ) {\n    total\n    data {\n      _id\n      trackingPlatform\n      trackingId\n      itemName\n      itemValueKIP\n      itemValueTHB\n      itemValueUSD\n      realItemValueKIP\n      realItemValueTHB\n      realItemValueUSD\n      receiverName\n      receiverPhone\n      description\n      isSummary\n      destSendDate\n      charge_on_shop\n      itemStatus\n      contactStatus\n      originSendDate\n      width\n      weight\n      isCod\n      isExtraItem\n      packagePrice\n      originReceiveDate\n      destReceiveDate\n      sendCompleteDate\n      isBackward\n      billNumber\n      originProvinceId {\n        provinceName\n      }\n      destProvinceId {\n        provinceName\n      }\n      originBranchId {\n        branch_name\n      }\n      destBranchId {\n        branch_name\n        branch_address\n        contactInfo\n      }\n      customerId {\n        id_list\n        full_name\n        contact_info\n      }\n      createdBy {\n        first_name\n        phone_number\n      }\n      originReceiveBy {\n        first_name\n        phone_number\n      }\n      providedBy {\n        _id\n      }\n    }\n  }\n}"
//   }' \
//   --compressed
