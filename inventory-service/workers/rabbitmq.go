package workers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type StockUpdateMessage struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

func StartStockConsumer(repo *repository.ProductRepository) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	queue, err := channel.QueueDeclare(
		"stock_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	go func() {
		for d := range msgs {
			var message StockUpdateMessage
			err := json.Unmarshal(d.Body, &message)

			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("Received stock update: ProductCode=%s, Quantity=%d", message.ProductCode, message.Quantity)
			currentProduct, err := repo.GetProductByCode(message.ProductCode)
			if err != nil {
				log.Printf("Failed to get product by code: %v", err)
				continue
			}

			newBalance := currentProduct.Balance - message.Quantity
			err = repo.UpdateProductBalance(message.ProductCode, newBalance)
			if err != nil {
				log.Printf("Failed to update product balance: %v", err)
				continue
			}

		}
	}()

	log.Printf("Stock consumer started. Waiting for messages...")

	<-sigchan
	log.Println("Stock consumer shut down.")

}
