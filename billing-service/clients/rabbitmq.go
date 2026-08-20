package clients

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type StockUpdateMessage struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

func PublishStockUpdate(productCode string, quantity int) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return err
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return err
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"stock_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return err
	}

	message := StockUpdateMessage{
		ProductCode: productCode,
		Quantity:    quantity,
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("Failed to publish a message:", err)
		return err
	}

	log.Printf("Published stock update for product %s with quantity %d", productCode, quantity)
	return nil
}
