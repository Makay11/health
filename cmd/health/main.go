package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Makay11/health"
)

func main() {
	config := getConfig()

	health.Start(health.Config{
		Port:           config.getString("port"),
		Path:           config.getString("path"),
		RequestTimeout: config.getDuration("requestTimeout"),
		CheckInterval:  config.getDuration("checkInterval"),
		Services:       config.getServices("services"),
	})
}

func getConfig() (config jsonObject) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	check(err)

	err = json.Unmarshal(data, &config)
	check(err)

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

type jsonObject map[string]any

func (object jsonObject) getArray(key string) []any {
	return getObjectValue[[]any](object, key)
}

func (object jsonObject) getDuration(key string) time.Duration {
	return time.Duration(object.getNumber(key)) * time.Millisecond
}

func (object jsonObject) getNumber(key string) float64 {
	return getObjectValue[float64](object, key)
}

func (object jsonObject) getString(key string) string {
	return getObjectValue[string](object, key)
}

func (object jsonObject) getServices(key string) (services []health.ServiceConfig) {
	rawServices := object.getArray(key)

	services = make([]health.ServiceConfig, len(rawServices))

	for i, rawService := range rawServices {
		o := jsonObject(castValue[map[string]any](rawService))

		services[i] = health.ServiceConfig{
			Name: o.getString("name"),
			Url:  o.getString("url"),
		}
	}

	return
}

func getObjectValue[T any](o jsonObject, key string) (value T) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Value of key %#v in JSON object has value %#v which cannot be cast to type %T.", key, o[key], value))
		}
	}()

	return castValue[T](o[key])
}

func castValue[T any](value any) T {
	v, ok := value.(T)

	if !ok {
		panic(fmt.Sprintf("Value %#v cannot be cast to type %T.", value, v))
	}

	return v
}
