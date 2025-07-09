package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/d4vz/rinha-de-backend-2025/pkg/store"
	"github.com/redis/go-redis/v9"
)

func EnqueuePayment(payment Payment) error {
	paymentJSON, err := json.Marshal(payment)

	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	err = store.GetRedis().LPush(context.Background(), "payments_queue", paymentJSON).Err()

	if err != nil {
		return fmt.Errorf("failed to enqueue payment: %w", err)
	}

	return nil
}

func DequeuePayment() (*Payment, error) {
	paymentJSON, err := store.GetRedis().RPop(context.Background(), "payments_queue").Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to dequeue payment: %w", err)
	}

	var paymentObj Payment

	if err := json.Unmarshal([]byte(paymentJSON), &paymentObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment: %w", err)
	}

	return &paymentObj, nil
}
