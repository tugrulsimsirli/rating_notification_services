package services

type RabbitMQServiceInterface interface {
	Publish(message string) error
}
