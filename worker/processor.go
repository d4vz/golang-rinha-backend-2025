package worker

import (
	"github.com/d4vz/rinha-de-backend-2025/payment"
	"github.com/gofiber/fiber/v2/log"
)

type WorkerPool struct {
	Workers        int
	Queue          chan *payment.Payment
	paymentService *payment.PaymentService
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		Workers:        workers,
		Queue:          make(chan *payment.Payment),
		paymentService: payment.NewPaymentService(),
	}
}

func (wp *WorkerPool) StartConsumer() {
	go func() {
		for {
			payment, err := payment.DequeuePayment()

			if err != nil {
				log.Errorf("Error dequeuing payment: %v", err)
				continue
			}

			if payment == nil {
				continue
			}

			wp.Queue <- payment
		}
	}()
}

func (wp *WorkerPool) StartProcessor() {
	for i := 0; i < wp.Workers; i++ {
		go wp.processPayments(wp.Queue)
	}
}

func (wp *WorkerPool) processPayments(q chan *payment.Payment) {
	for p := range q {
		err := wp.paymentService.SendToDefaultProcessor(p)

		if err != nil {
			log.Errorf("Erro ao enviar pagamento para o processador fallback: %v. Reenfileirando na fila persistente...", err)
			payment.EnqueuePayment(*p)
			continue
		}

		// if err != nil {
		// 	log.Errorf("Erro ao enviar pagamento para o processador default: %v. Tentando fallback...", err)
		// 	p.Processor = config.FallbackProcessorName
		// 	err = wp.paymentService.SendToFallbackProcessor(p)

		// }

		if err := payment.CreatePayment(*p); err != nil {
			log.Errorf("Erro ao criar registro do pagamento: %v. Pagamento será descartado.", err)
		}
	}
}
