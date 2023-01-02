#!/bin/sh
echo "Copy environment file"
yes | cp -rf build/weather-api-env /root/go/env/weather-api-env
echo "Build go application"
GOOS=linux GOARCH=amd64 go build -o weather-go-api main.go
echo "Restart service"
systemctl restart weather-go-api
systemctl status weather-go-api