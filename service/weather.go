package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather/common/log"
	"weather/common/model"
	"weather/common/response"

	"github.com/go-resty/resty/v2"
)

type (
	IWeather interface {
		GetWeatherData(ctx context.Context) (int, interface{})
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

func (service *Weather) GetWeatherData(ctx context.Context) (int, interface{}) {
	resp := []model.WeatherData{}
	cities := []string{"Ho Chi Minh", "Ha Noi", "Da Nang"}
	for _, val := range cities {
		queryParams := map[string]string{}
		queryParams["q"] = val
		queryParams["limit"] = "5"
		queryParams["appid"] = service.config.Key
		client := resty.New()
		client.SetTimeout(time.Second * 3)
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
		url := fmt.Sprintf("%s/geo/1.0/direct", service.config.Domain)
		res, err := client.R().
			SetQueryParams(queryParams).
			ForceContentType("application/json").
			Get(url)
		if err != nil {
			log.Error(err)
			continue
		}
		result := []model.WeatherCoordinateLocation{}
		if err := json.Unmarshal(res.Body(), &result); err != nil {
			log.Error(err)
			continue
		}
		if len(result) > 0 {
			queryParams := map[string]string{}
			queryParams["lat"] = fmt.Sprintf("%f", result[0].Lat)
			queryParams["lon"] = fmt.Sprintf("%f", result[0].Lon)
			queryParams["appid"] = service.config.Key
			client := resty.New()
			client.SetTimeout(time.Second * 3)
			client.SetTLSClientConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
			url := fmt.Sprintf("%s/data/2.5/weather", service.config.Domain)
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
