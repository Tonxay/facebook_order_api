package cons

var Ordered = "ordered"
var OrderCancelled = "order_cancelled"
var PaymentCompleted = "payment_completed"

var OrderStatus = []string{
	"ordered",                // 0: ສັ່ງຊື້ (Order placed)
	"waiting_to_pack",        // 1: ລໍຖ້າແພັກເຄື່ອງ (Waiting to pack)
	"packed",                 // 2: ແພັກເຄື່ອງແລ້ວ (Packed)
	"shipped",                // 3: ສົ່ງແລ້ວ (Shipped)
	"customer_bill_notified", // 4: ແຈ້ງບິນໃຫ້ລູກຄ້າແລ້ວ (Customer bill notified)
	"delivery_complete",      // 5: ສົ່ງສຳເລັດ (Delivery complete)
	"payment_completed",      // 6: ສຳລະເງິນແລ້ວ (Payment completed)
	"order_cancelled",        // 7: ຍົກເລີກອໍເດີ (Order cancelled)
	"return_to_sender",       // 8: ພັດສະດຸຕີກັບ (Returned to sender)
	"customer_notified",      // 9: ແຈ້ງລູກແລ້ວ (Customer notified)
}

var OrderStatusLao = map[string]string{
	"ordered":                "ສັ່ງຊື້",
	"waiting_to_pack":        "ລໍຖ້າແພັກເຄື່ອງ",
	"packed":                 "ແພັກເຄື່ອງແລ້ວ",
	"shipped":                "ສົ່ງແລ້ວ",
	"customer_bill_notified": "ແຈ້ງບິນໃຫ້ລູກຄ້າແລ້ວ",
	"delivery_complete":      "ສົ່ງສຳເລັດ",
	"payment_completed":      "ສຳລະເງິນແລ້ວ",
	"order_cancelled":        "ຍົກເລີກອໍເດີ",
	"return_to_sender":       "ພັດສະດຸຕີກັບ",
	"customer_notified":      "ແຈ້ງລູກແລ້ວ",
}

// OrderStatusTransitions defines allowed next steps from each status
var OrderStatusTransitions = map[string]string{
	"ordered":                "ordered",
	"waiting_to_pack":        "ordered",
	"packed":                 "waiting_to_pack",
	"shipped":                "packed",
	"customer_bill_notified": "shipped",
	"delivery_complete":      "customer_bill_notified",
	"customer_notified":      "delivery_complete",
	"payment_completed":      "customer_notified",
	"return_to_sender":       "customer_notified",
	"order_cancelled":        "order_cancelled",
}
