package main

import (
	"log"
	"notification_service/internal/app/services"
	"notification_service/internal/http/handlers"

	_ "notification_service/docs" // Import the generated docs

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tugrulsimsirli/rabbitmq"
)

func main() {
	// Initialize Echo
	e := echo.New()

	rabbitMQService, err := rabbitmq.NewRabbitMQService("amqp://guest:guest@rabbitmq:5672/", "notification_queue")
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQService.Close()

	notificationService := services.NotificationService{RabbitMQService: rabbitMQService}
	notificationHandler := handlers.NotificationHandler{NotificationService: notificationService}

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	e.GET("/notifications", notificationHandler.GetNotifications)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
