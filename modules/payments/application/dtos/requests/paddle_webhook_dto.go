package paymentsrequestsdtos

type PaddleWebhookPayload struct {
	EventID   string            `json:"event_id"`
	EventType string            `json:"event_type"`
	Data      PaddleWebhookData `json:"data"`
}

type PaddleWebhookData struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	CustomerID   *string                `json:"customer_id"`
	CurrencyCode string                 `json:"currency_code"`
	Total        string                 `json:"total"`
	Items        []PaddleWebhookItem    `json:"items"`
	CustomData   map[string]interface{} `json:"custom_data"`
}

type PaddleWebhookItem struct {
	Price    PaddleWebhookPrice `json:"price"`
	Quantity int                `json:"quantity"`
}

type PaddleWebhookPrice struct {
	ID string `json:"id"`
}
