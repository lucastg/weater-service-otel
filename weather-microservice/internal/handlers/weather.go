package handlers

import (
	"net/http"
	"weather-microservice/internal/services"
	"weather-microservice/internal/utils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func GetWeatherByZipcode(c *gin.Context) {
	tracer := otel.Tracer("weather-service")
	ctx, span := tracer.Start(c, "ClimaHandler")
	defer span.End()

	zipcode := c.Param("zipcode")

	if len(zipcode) != 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	ctx, spanLocation := tracer.Start(ctx, "BuscarLocalizacaoPorCEP")
	location, err := services.GetLocationByZipcode(zipcode)
	spanLocation.End()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "can not find zipcode"})
		return
	}

	ctx, spanWeather := tracer.Start(ctx, "BuscarClima")
	tempC, err := services.GetTemperature(location.City)
	spanWeather.End()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch temperature"})
		return
	}

	tempF := utils.CelsiusToFahrenheit(tempC)
	tempK := utils.CelsiusToKelvin(tempC)

	c.JSON(http.StatusOK, gin.H{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	})
}
