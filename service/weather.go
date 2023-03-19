package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather/common/log"
	"weather/common/model"
	"weather/common/response"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type (
	IWeather interface {
		GetWeatherData(ctx *gin.Context) (int, interface{})
		GetHistorical(ctx *gin.Context) (int, interface{})
	}
	Weather struct {
		config WeatherConfig
	}
	WeatherConfig struct {
		Domain string
		Key    string
	}
)

var weather IWeather

func NewWeather(config WeatherConfig) IWeather {
	return &Weather{
		config: config,
	}
}

func (weather *Weather) GetLocation(location string) []model.WeatherCoordinateLocation {
	result := []model.WeatherCoordinateLocation{}
	queryParams := map[string]string{}
	queryParams["q"] = location
	queryParams["limit"] = "5"
	queryParams["appid"] = weather.config.Key
	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	url := fmt.Sprintf("%s/geo/1.0/direct", weather.config.Domain)
	res, err := client.R().
		SetQueryParams(queryParams).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		log.Error(err)
	}
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		log.Error(err)
	}
	return result
}

func (weather *Weather) GetWeatherData(ctx *gin.Context) (int, interface{}) {
	resp := []model.WeatherData{}
	cities := []string{"Ho Chi Minh", "Ha Noi", "Da Nang"}
	for _, val := range cities {
		location := weather.GetLocation(val)
		if len(location) > 0 {
			queryParams := map[string]string{}
			queryParams["lat"] = fmt.Sprintf("%f", location[0].Lat)
			queryParams["lon"] = fmt.Sprintf("%f", location[0].Lon)
			queryParams["appid"] = weather.config.Key
			client := resty.New()
			client.SetTimeout(time.Second * 3)
			client.SetTLSClientConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
			url := fmt.Sprintf("%s/data/2.5/weather", weather.config.Domain)
			res, err := client.R().
				SetHeader("content-type", "application/json").
				SetQueryParams(queryParams).
				ForceContentType("application/json").
				Get(url)
			if err != nil {
				log.Error(err)
				continue
			}
			result := model.WeatherData{}
			if err := json.Unmarshal(res.Body(), &result); err != nil {
				log.Error(err)
				continue
			}
			result.NameLocation = val
			resp = append(resp, result)
		}
	}

	return response.Data(http.StatusOK, resp)
}

func (weather *Weather) GetHistorical(c *gin.Context) (int, interface{}) {
	location := weather.GetLocation("Ho Chi Minh")

	queryParams := map[string]string{}
	queryParams["lat"] = fmt.Sprintf("%f", location[0].Lat)
	queryParams["lon"] = fmt.Sprintf("%f", location[0].Lon)
	queryParams["appid"] = weather.config.Key
	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	url := fmt.Sprintf("%s/data/2.5/forecast", weather.config.Domain)
	res, err := client.R().
		SetHeader("content-type", "application/json").
		SetQueryParams(queryParams).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	result := model.Historical{}
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Data(http.StatusOK, result)
}
