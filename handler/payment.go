package handler

import (
	"fmt"
	"time"

	"github.com/d4vz/rinha-de-backend-2025/payment"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Summary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentSummary struct {
	Default  Summary `json:"default"`
	Fallback Summary `json:"fallback"`
}

func ProcessPayment(c *fiber.Ctx) error {
	var p payment.Payment

	if err := c.BodyParser(&p); err != nil {
		log.Errorf("failed to parse payment: %v", err)
		return c.Status(fiber.StatusBadRequest).Next()
	}

	payment.EnqueuePayment(p)

	return c.SendStatus(fiber.StatusAccepted)
}

func parseDate(dateStr string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}

	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02",
		time.RFC3339Nano,
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}
func GetPaymentSummary(c *fiber.Ctx) error {
	from, err := parseDate(c.Query("from"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Next()
	}

	to, err := parseDate(c.Query("to"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Next()
	}

	summary, err := payment.GetPaymentsSummary(from, to)

	if err != nil {
		log.Errorf("failed to get summary: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(summary)
}
