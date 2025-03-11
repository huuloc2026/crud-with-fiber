package config

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewRabbitMQConnection(cfg *RabbitMQConfig) (*amqp.Connection, *amqp.Channel) {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare the notifications queue
	_, err = ch.QueueDeclare(
		"notifications", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	return conn, ch
}
