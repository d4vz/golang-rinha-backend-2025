package handler

import (
	"github.com/d4vz/rinha-de-backend-2025/payment"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func PurgePayments(c *fiber.Ctx) error {
	if err := payment.PurgePayments(); err != nil {
		log.Errorf("failed to purge payments: %v", err)
		return c.Status(fiber.StatusInternalServerError).Next()
	}

	return c.SendStatus(fiber.StatusOK)
}
