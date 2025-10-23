package events

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(eventType string, payload interface{}) error
}

type RabbitMQPublisher struct {
	channel *amqp.Channel
	conn    *amqp.Connection
}

func (p *RabbitMQPublisher) Close() {
	p.channel.Close()
	p.conn.Close()
}

func NewRabbitMQPublisher() (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		"order_events", // name
		"topic",        // type
		true,           // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	return &RabbitMQPublisher{channel: ch, conn: conn}, nil
}

func (p *RabbitMQPublisher) Publish(eventType string, payload interface{}) error {

	body, err := json.Marshal(payload)
	log.Println("📤 OrderCreatedEvent sending to RabbitMQ:", string(body))
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"order_events",
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
