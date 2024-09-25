package main

import (
	"fmt"
	"log"
	"notification_service/config"
	"notification_service/internal/app/services"
	"notification_service/internal/http/handlers"

	_ "notification_service/docs" // Import the generated docs

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tugrulsimsirli/rabbitmq"
)

func main() {
	// Config dosyasını yükle
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Echo
	e := echo.New()

	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQService.Close()

	notificationService := services.NotificationService{
		RabbitMQService: rabbitMQService,
		Config:          cfg,
	}
	notificationHandler := handlers.NotificationHandler{NotificationService: notificationService}

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	e.GET("/notifications", notificationHandler.GetNotifications)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}
