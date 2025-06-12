package request

type StatusOrderRequest struct {
	Statuses  []string `json:"statuses"`
	IsCancell bool     `json:"is_cancel"`
}
