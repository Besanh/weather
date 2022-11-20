package model

import "encoding/json"

type WeatherCoordinateLocation struct {
	Name       string          `json:"name"`
	LocalNames json.RawMessage `json:"local_names"`
	Lat        float64         `json:"lat"`
	Lon        float64         `json:"lon"`
	Country    string          `json:"country"`
}

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []Weather `json:"weather"`
	Base    string    `json:"base"`
	Main    struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  float64 `json:"pressure"`
		Humidity  float64 `json:"humidity"`
	} `json:"main"`
	Visibility float64 `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All float64 `json:"all"`
	} `json:"clouds"`
	Dt  uint64 `json:"dt"`
	Sys struct {
		Type    int64  `json:"type"`
		Id      int64  `json:"id"`
		Country string `json:"country"`
		Sunrise uint64 `json:"sunrise"`
		Sunset  uint64 `json:"sunset"`
	} `json:"sys"`
	Timezone int64  `json:"timezone"`
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type Weather struct {
	Id          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
}
