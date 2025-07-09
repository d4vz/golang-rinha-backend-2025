package payment

import "time"

type Payment struct {
	ID            int       `json:"-"`
	CorrelationID string    `json:"correlationId"`
	Amount        float64   `json:"amount"`
	Processor     string    `json:"processor,omitempty"`
	RequestedAt   time.Time `json:"requestedAt,omitempty"`
	RetryCount    int       `json:"-"`
}

type Summary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentsSummary struct {
	Default  Summary `json:"default"`
	Fallback Summary `json:"fallback"`
}
