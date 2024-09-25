package services

import "github.com/streadway/amqp"

// RabbitMQServiceInterface defines methods for RabbitMQ operations
type RabbitMQServiceInterface interface {
	CreateChannel(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error)
	CloseChannel(ch *amqp.Channel) error
}
