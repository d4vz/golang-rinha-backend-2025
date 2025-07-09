package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/d4vz/rinha-de-backend-2025/config"
	"github.com/d4vz/rinha-de-backend-2025/pkg/store"
)

func CreatePayment(payment Payment) error {
	db := store.GetDB()

	_, err := db.Exec(context.Background(),
		"INSERT INTO payments (correlation_id, amount, processor, processed_at) VALUES ($1, $2, $3, $4)",
		payment.CorrelationID, payment.Amount, payment.Processor, payment.RequestedAt)

	return err
}

func GetPaymentsSummary(from, to time.Time) (PaymentsSummary, error) {
	db := store.GetDB()

	rows, err := db.Query(context.Background(),
		"SELECT processor, COUNT(*), SUM(amount) FROM payments WHERE processed_at BETWEEN $1 AND $2 GROUP BY processor", from, to)

	if err != nil {
		return PaymentsSummary{}, err
	}

	defer rows.Close()

	summary := PaymentsSummary{}

	for rows.Next() {
		var processor string
		var count int
		var amount float64

		if err := rows.Scan(&processor, &count, &amount); err != nil {
			return PaymentsSummary{}, err
		}

		switch processor {
		case config.DefaultProcessorName:
			summary.Default.TotalRequests += count
			summary.Default.TotalAmount += amount
		case config.FallbackProcessorName:
			summary.Fallback.TotalRequests += count
			summary.Fallback.TotalAmount += amount
		default:
			return PaymentsSummary{}, fmt.Errorf("invalid processor: %s", processor)
		}
	}

	return summary, nil
}

func PurgePayments() error {
	db := store.GetDB()

	_, err := db.Exec(context.Background(), "DELETE FROM payments")

	return err
}
