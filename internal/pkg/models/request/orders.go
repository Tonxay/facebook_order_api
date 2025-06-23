package request

type StatusOrderRequest struct {
	Statuses []string `json:"statuses"`
	IsCancel bool     `json:"is_cancel"`
	Tel      string   `json:"tel"`
}
