package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"weather/api"
	apiV1 "weather/api/v1"
	"weather/service"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Dir           string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port          string
	LogType       string
	LogLevel      string
	LogFile       string
	Storage       string
	Auth          string
	DB            string
	ElasticSearch string
	Redis         string
	RabbitMq      string
	EventSocket   string
	Callback      string
	Sftp          string
}

var config Config

func init() {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Error(err)
	}
	time.Local = loc
	if err := env.Parse(&config); err != nil {
		log.Error("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Dir:           config.Dir,
		Port:          viper.GetString("main.port"),
		LogType:       viper.GetString("main.log_type"),
		LogFile:       viper.GetString("main.log_file"),
		LogLevel:      viper.GetString("main.log_level"),
		Storage:       viper.GetString("main.storage"),
		Auth:          viper.GetString("main.auth"),
		DB:            viper.GetString("main.db"),
		ElasticSearch: viper.GetString("main.elastic_search"),
		Redis:         viper.GetString("main.redis"),
		RabbitMq:      viper.GetString("main.rabbitmq"),
		EventSocket:   viper.GetString("main.event_socket"),
		Callback:      viper.GetString("main.callback"),
		Sftp:          viper.GetString("main.sftp"),
	}
	config = cfg
}

func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	if err := createNewLogFile(config.LogFile); err != nil {
		log.Error(err)
	}
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)

	s := api.NewServer()
	weatherService := service.NewWeather(service.WeatherConfig{
		Domain: viper.GetString("open_weather.domain"),
		Key:    viper.GetString("open_weather.key"),
	})
	apiV1.NewApiWeather(s.Engine, viper.GetString("open_weather.signature"), weatherService)

	s.Start(config.Port)
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.Error("main ", "Log File "+cfg.LogFile+" error")
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}

func createNewLogFile(logDir string) error {
	files, err := os.ReadDir("tmp")
	if err != nil {
		return err
	}
	last10dayUnix := time.Now().Add(-1 * 24 * time.Hour).Unix()
	for _, f := range files {
		tmp := strings.Split(f.Name(), ".")
		if len(tmp) > 2 {
			fileUnix, err := strconv.Atoi(tmp[2])
			if err != nil {
				return err
			} else if int64(fileUnix) < last10dayUnix {
				if err := os.Remove("tmp/" + f.Name()); err != nil {
					return err
				}
			}
		}
	}
	_, err = os.Stat(logDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err := os.Rename(logDir, fmt.Sprintf(logDir+".%d", time.Now().Unix())); err != nil {
		return err
	}
	return nil
}
