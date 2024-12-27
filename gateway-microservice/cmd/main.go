package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	services "gateway-microservice/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
)

type ZipcodeRequest struct {
	CEP string `json:"cep"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	loadEnv()

	tracerProvider := services.InitTracer(os.Getenv("ZIPKIN_URL"), "weather-gateway")
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.POST("/weather", handleZipcodeRequest)

	log.Println("Gateway running on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func loadEnv() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	godotenv.Load(filepath.Join(dir, ".env"))
}

func validateZipcode(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func handleZipcodeRequest(ctx *gin.Context) {
	var request ZipcodeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if !validateZipcode(request.CEP) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid zipcode"})
		return
	}

	tracer := otel.Tracer("gateway-service")
	otelCtx, span := tracer.Start(ctx.Request.Context(), "Call Weather Service")
	defer span.End()

	resp, err := callWeatherService(otelCtx, request.CEP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error calling Weather Service"})
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading Weather Service response"})
		return
	}

	ctx.Data(resp.StatusCode, "application/json", responseBody)
}

func callWeatherService(ctx context.Context, cep string) (*http.Response, error) {
	tracer := otel.Tracer("gateway-service")
	_, span := tracer.Start(ctx, "Weather Service Request")
	defer span.End()

	weatherServiceUrl := os.Getenv("WEATHER_API_URL")
	url := fmt.Sprintf("%s/weather/%s", weatherServiceUrl, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
