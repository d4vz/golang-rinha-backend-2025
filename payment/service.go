package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/d4vz/rinha-de-backend-2025/config"
	"github.com/sony/gobreaker"
)

type PaymentService struct {
	DefaultProcessorCB  *gobreaker.CircuitBreaker
	FallbackProcessorCB *gobreaker.CircuitBreaker
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		DefaultProcessorCB: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "DefaultProcessor",
			MaxRequests: 10,
			Interval:    1 * time.Second,
			Timeout:     1 * time.Second, // Tempo que o circuito permanece aberto após falha
		}),
		FallbackProcessorCB: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "FallbackProcessor",
			MaxRequests: 10,
			Interval:    1 * time.Second,
			Timeout:     1 * time.Second, // Tempo que o circuito permanece aberto após falha
		}),
	}
}

// Os timeouts abaixo estão de acordo com o rinha.js:
// - DefaultProcessor: até 2 segundos de timeout
// - FallbackProcessor: até 5 segundos de timeout
func (wp *PaymentService) SendToDefaultProcessor(p *Payment) error {
	_, err := wp.DefaultProcessorCB.Execute(func() (interface{}, error) {
		p.Processor = config.DefaultProcessorName
		p.RequestedAt = time.Now()
		httpClient := &http.Client{Timeout: 2 * time.Second} // conforme rinha.js
		url := config.GetEnvOrDefault("PROCESSOR_DEFAULT_URL", "http://localhost:8001") + "/payments"
		jsonPayload, err := json.Marshal(p)

		if err != nil {
			return nil, err
		}

		response, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonPayload))

		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("default processor returned non-200 status: %d", response.StatusCode)
		}

		return nil, nil
	})

	return err
}

func (wp *PaymentService) SendToFallbackProcessor(p *Payment) error {
	_, err := wp.FallbackProcessorCB.Execute(func() (interface{}, error) {
		p.Processor = config.FallbackProcessorName
		p.RequestedAt = time.Now()
		httpClient := &http.Client{Timeout: 5 * time.Second} // conforme rinha.js
		url := config.GetEnvOrDefault("PROCESSOR_FALLBACK_URL", "http://localhost:8002") + "/payments"
		jsonPayload, err := json.Marshal(p)

		if err != nil {
			return nil, err
		}

		response, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonPayload))

		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fallback processor returned non-200 status: %d", response.StatusCode)
		}

		return nil, nil
	})

	return err
}
