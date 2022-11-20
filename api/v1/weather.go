package v1

import (
	"net/http"
	"weather/service"

	"github.com/gin-gonic/gin"
)

type Weather struct {
	weatherService service.IWeather
}

func NewApiWeather(c *gin.Engine, signature string, weatherService service.IWeather) {
	handler := Weather{
		weatherService: weatherService,
	}

	Group := c.Group("v1/weather")
	{
		Group.GET("data", ValidHeader(signature), handler.GetWeatherData)
	}
}

func ValidHeader(signature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("open-weather-signature")
		if token != signature {
			c.JSON(http.StatusUnauthorized,
				map[string]interface{}{
					"error": http.StatusText(http.StatusUnauthorized),
				})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (handler *Weather) GetWeatherData(c *gin.Context) {
	code, result := handler.weatherService.GetWeatherData(c)
	c.JSON(code, result)
}