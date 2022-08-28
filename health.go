package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Makay11/health/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Port           string
	Path           string
	RequestTimeout time.Duration
	CheckInterval  time.Duration
	Services       []ServiceConfig
}

type ServiceConfig struct {
	Name string
	Url  string
}

func Start(config Config) error {
	httpClient := http.Client{
		Timeout: config.RequestTimeout,
	}

	services := make([]service, len(config.Services))

	for i, serviceConfig := range config.Services {
		services[i] = service{
			name: serviceConfig.Name,
			url:  serviceConfig.Url,
		}

		go services[i].check(&httpClient, config.CheckInterval)
	}

	app := fiber.New()

	app.Get(config.Path, func(c *fiber.Ctx) error {
		return c.JSON(services)
	})

	return app.Listen(config.Port)
}

type service struct {
	name              string
	url               string
	ok                bool
	responseTimeMilli int64
	responseTimeMicro int64
	lastSeen          time.Time
	error             error
}

func (service *service) MarshalJSON() ([]byte, error) {
	m := fiber.Map{
		"name":              service.name,
		"url":               service.url,
		"ok":                service.ok,
		"responseTimeMilli": service.responseTimeMilli,
		"responseTimeMicro": service.responseTimeMicro,
	}

	if service.lastSeen.IsZero() {
		m["lastSeen"] = nil
	} else {
		m["lastSeen"] = service.lastSeen
	}

	if service.error == nil {
		m["error"] = nil
	} else {
		m["error"] = service.error.Error()
	}

	return json.Marshal(m)
}

func (service *service) check(httpClient *http.Client, checkInterval time.Duration) {
	for {
		startTime := time.Now()
		resp, err := httpClient.Head(service.url)
		endTime := time.Now()

		if err != nil {
			utils.Logger.Println(err)

			service.ok = false
			service.responseTimeMilli = 0
			service.responseTimeMicro = 0
			service.error = err
		} else {
			if resp.StatusCode != 200 {
				err := fmt.Errorf("%v returned invalid status code %v", service.url, resp.StatusCode)
				utils.Logger.Println(err)

				service.ok = false
				service.error = err
			} else {
				// logger.Printf("%v returned valid status code %v", service.Url, resp.StatusCode)

				service.ok = true
				service.error = nil
			}

			service.responseTimeMilli = endTime.UnixMilli() - startTime.UnixMilli()
			service.responseTimeMicro = endTime.UnixMicro() - startTime.UnixMicro()
			service.lastSeen = endTime
		}

		time.Sleep(checkInterval)
	}
}
