package main

import (
	"context"
	"log"
	"os"

	"weather-microservice/internal/handlers"
	"weather-microservice/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	tracerProvider := services.InitTracer(os.Getenv("ZIPKIN_URL"), "weather-gateway")
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WeatherAPI key not set")
	}

	r := gin.Default()

	r.GET("/weather/:zipcode", handlers.GetWeatherByZipcode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Weater service running on port " + port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}
}
