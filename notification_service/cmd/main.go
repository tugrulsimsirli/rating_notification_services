package main

import (
	"database/sql"
	"log"
	"notification_service/internal/app/services"
	"notification_service/internal/db/repositories"
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

	// Initialize DB connection
	db, err := sql.Open("postgres", "host=db user=rating password=rating dbname=ratingdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rabbitMQService, err := rabbitmq.NewRabbitMQService("amqp://guest:guest@rabbitmq:5672/", "notification_queue")
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQService.Close()

	notificationRepo := repositories.NotificationRepository{DB: db}
	notificationService := services.NotificationService{NotificationRepository: notificationRepo}
	notificationHandler := handlers.NotificationHandler{
		NotificationService: notificationService,
		RabbitMQService:     rabbitMQService,
	}

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	e.GET("/notifications/:providerID", notificationHandler.GetNotifications)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}
