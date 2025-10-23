package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/repository"
	"github.com/streadway/amqp"
)

type OrderStatusConsumer struct {
	ch        *amqp.Channel
	orderRepo repository.OrderRepository
}

const (
	ExchangeName = "order_events"
	ExchangeType = "topic"
)

func NewOrderStatusConsumer(orderRepo repository.OrderRepository) (*OrderStatusConsumer, func(), error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@rabbitmq:5672/"
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, func() {}, fmt.Errorf("dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, func() {}, fmt.Errorf("ch: %w", err)
	}

	// Ensure exchange exists
	if err := ch.ExchangeDeclare(ExchangeName, ExchangeType, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, func() {}, err
	}

	// Declare & bind queues for status events
	for _, rk := range []string{"order.reserved", "order.rejected"} {
		q, err := ch.QueueDeclare(rk, true, false, false, false, nil)
		if err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return nil, func() {}, err
		}
		if err := ch.QueueBind(q.Name, rk, ExchangeName, false, nil); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return nil, func() {}, err
		}
	}

	cleanup := func() { _ = ch.Close(); _ = conn.Close() }
	return &OrderStatusConsumer{ch: ch, orderRepo: orderRepo}, cleanup, nil
}

func (c *OrderStatusConsumer) Start() error {
	consume := func(rk string, handler func([]byte) error) error {
		msgs, err := c.ch.Consume(rk, "", false, false, false, false, nil)
		if err != nil {
			return err
		}
		go func() {
			for msg := range msgs {
				if err := handler(msg.Body); err != nil {
					log.Printf("handler error (%s): %v", rk, err)
					_ = msg.Nack(false, true)
					continue
				}
				_ = msg.Ack(false)
			}
		}()
		return nil
	}

	if err := consume("order.reserved", func(b []byte) error {
		var ev domain.OrderReservedEvent
		if err := json.Unmarshal(b, &ev); err != nil {
			return err
		}
		return c.orderRepo.UpdateOrderStatus(ev.OrderID, domain.Completed)
	}); err != nil {
		return err
	}

	if err := consume("order.rejected", func(b []byte) error {
		var ev domain.OrderRejectedEvent
		if err := json.Unmarshal(b, &ev); err != nil {
			return err
		}
		return c.orderRepo.UpdateOrderStatus(ev.OrderID, domain.Cancelled)
	}); err != nil {
		return err
	}

	return nil
}
