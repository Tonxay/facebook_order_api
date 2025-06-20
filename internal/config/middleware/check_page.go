package middleware

import "os"

func CheckPageId(sendId string, recipientID string) (string, string) {
	pageId := os.Getenv("PAGE_ID")
	pageToken := os.Getenv("PAGE_ACCESS_TOKEN")

	pageNanaId := os.Getenv("PAGE_NANA_ID")
	pageNanaToken := os.Getenv("PAGE_ACCESS_TOKEN_NANA")
	pageNanaShopId := os.Getenv("PAGE_NANA_SHOP_ID")
	pageNanaShopToken := os.Getenv("PAGE_ACCESS_TOKEN_NANA_SHOP")

	// dd
	if sendId == pageId || recipientID == pageId {
		return pageId, pageToken
		// nana
	} else if sendId == pageNanaId || recipientID == pageNanaId {
		return pageNanaId, pageNanaToken
		// na shop
	} else if sendId == pageNanaShopId || recipientID == pageNanaShopToken {
		return pageNanaShopId, pageNanaShopToken
		// dd
	} else {
		return pageId, pageToken
	}

}
