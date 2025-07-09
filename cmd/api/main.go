package main

import (
	"fmt"
	"log"

	"github.com/d4vz/rinha-de-backend-2025/config"
	"github.com/d4vz/rinha-de-backend-2025/handler"
	"github.com/d4vz/rinha-de-backend-2025/pkg/store"
	"github.com/d4vz/rinha-de-backend-2025/worker"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := store.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	if err := store.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := store.MigrateDB(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := fiber.New(fiber.Config{
		DisableKeepalive: true,
	})

	workerPool := worker.NewWorkerPool(50)

	for i := 0; i < 4; i++ {
		go workerPool.StartConsumer()
	}

	go workerPool.StartProcessor()

	app.Get("/payments-summary", handler.GetPaymentSummary)
	app.Post("/payments", handler.ProcessPayment)
	app.Post("/purge-payments", handler.PurgePayments)

	port := config.GetEnvOrDefault("PORT", "8080")
	log.Printf("Starting server on port %s ...", port)
	app.Listen(fmt.Sprintf(":%s", port))
}
