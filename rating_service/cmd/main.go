package main

import (
	"database/sql"
	"fmt"
	"log"
	"rating_service/config"
	"rating_service/internal/app/services"
	"rating_service/internal/db/repositories"
	"rating_service/internal/http/handlers"

	_ "rating_service/docs" // Import the generated docs

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tugrulsimsirli/rabbitmq"
)

func main() {
	// Config dosyasından ayarları yükle
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Echo
	e := echo.New()

	// Initialize DB connection
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQService.Close()

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Initialize repositories, services, and handlers
	ratingRepo := repositories.RatingRepository{DB: db}
	ratingService := services.RatingService{
		RatingRepository: ratingRepo,
		RabbitMQService:  rabbitMQService,
	}
	ratingHandler := handlers.RatingHandler{
		RatingService: ratingService,
	}

	// Routes
	e.POST("/submit-rating", ratingHandler.SubmitRating)
	e.GET("/average-rating/:providerID", ratingHandler.GetAverageRating)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}
