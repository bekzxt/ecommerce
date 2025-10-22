package events

import (
	"encoding/json"
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"log"
	"os"

	"github.com/bekzxt/e-commerce/inventory-service/internal/interfaces/usecase"
	"github.com/streadway/amqp"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (c *RabbitMQConsumer) Channel() *amqp.Channel {
	return c.channel
}

type OrderCreatedConsumer struct {
	channel     *amqp.Channel
	inventoryUC usecase.Inventory
	publisher   *RabbitMQPublisher
}

func NewRabbitMQConsumer() (*RabbitMQConsumer, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{conn: conn, channel: ch}, nil
}

func (c *RabbitMQConsumer) Close() {
	_ = c.channel.Close()
	_ = c.conn.Close()
}

func NewOrderCreatedConsumer(channel *amqp.Channel, inventoryUC usecase.Inventory, publisher *RabbitMQPublisher) *OrderCreatedConsumer {
	return &OrderCreatedConsumer{channel: channel, inventoryUC: inventoryUC, publisher: publisher}
}

func (c *OrderCreatedConsumer) Consume() error {
	// Объявляем очередь (если нет)
	q, err := c.channel.QueueDeclare("order.created", true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := c.channel.QueueBind(q.Name, "order.created", ExchangeName, false, nil); err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		"order.created",
		"",
		false, // manual Ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Горутинный обработчик сообщений
	go func() {
		for msg := range msgs {
			type OrderCreated struct {
				OrderID string                `json:"order_id"` // UUID
				UserID  string                `json:"user_id"`
				Total   float64               `json:"total"`
				Status  string                `json:"status"`
				Items   []domain.OrderItemInv `json:"items"`
			}
			var event OrderCreated
			// Вызываем юзкейс бизнес-логики
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("❌ Failed to parse message:", err)
				_ = msg.Nack(false, false)
				continue
			}

			log.Printf("📦 Received order: %d", event.OrderID)
			for _, item := range event.Items {
				log.Printf("🔥 Parsed item from RabbitMQ: product_id=%v, quantity=%v", item.ProductID, item.Quantity)
			}
			ok, err := c.inventoryUC.ReserveItems(event.Items)
			if err != nil {
				log.Println("❌ ReserveItems failed:", err)
				if c.publisher != nil {
					_ = c.publisher.Publish("order.rejected", map[string]string{"order_id": event.OrderID})
				}
				_ = msg.Nack(false, false)
				continue
			}

			if !ok {
				log.Println("⚠️ Not enough stock for order:", event.OrderID)
				_ = msg.Ack(false)
				if c.publisher != nil {
					_ = c.publisher.Publish("order.rejected", map[string]string{"order_id": event.OrderID})
				}
				continue
			}

			log.Println("✅ Reserved items for order:", event.OrderID)
			_ = msg.Ack(false)
			if c.publisher != nil {
				_ = c.publisher.Publish("order.reserved", map[string]string{"order_id": event.OrderID})
			}
		}
	}()

	return nil
}
