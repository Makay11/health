package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Makay11/health"
)

func main() {
	health.Start(getConfig())
}

func getConfig() (config health.Config) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	check(err)

	err = json.Unmarshal(data, &config)
	check(err)

	config.RequestTimeout *= time.Millisecond
	config.CheckInterval *= time.Millisecond

	return
}

func getConfigPath() (configPath string) {
	if len(os.Args) > 2 {
		panic("Too many arguments. Only 0 or 1 are supported.")
	}

	if len(os.Args) == 2 {
		configPath = os.Args[1]
	} else {
		configPath = os.Getenv("HEALTH_CONFIG")
	}

	if configPath == "" {
		panic("No config path provided. Provide it as argument or set HEALTH_CONFIG.")
	}

	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
