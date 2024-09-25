package main

import (
	"database/sql"
	"log"
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

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Initialize repositories, services, and handlers
	ratingRepo := repositories.RatingRepository{DB: db}
	ratingService := services.RatingService{RatingRepository: ratingRepo}
	ratingHandler := handlers.RatingHandler{
		RatingService:   ratingService,
		RabbitMQService: rabbitMQService,
	}

	// Routes
	e.POST("/submit-rating", ratingHandler.SubmitRating)
	e.GET("/average-rating/:providerID", ratingHandler.GetAverageRating)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
