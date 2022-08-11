package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Start(port, path string, requestTimeout, checkInterval time.Duration, services []Service) error {
	httpClient := http.Client{
		Timeout: requestTimeout,
	}

	for i := range services {
		go services[i].check(&httpClient, checkInterval)
	}

	app := fiber.New()

	app.Get(path, func(c *fiber.Ctx) error {
		return c.JSON(services)
	})

	return app.Listen(port)
}

type Service struct {
	Name              string
	Url               string
	ok                bool
	responseTimeMilli int64
	responseTimeMicro int64
	lastSeen          time.Time
	error             error
}

func (service *Service) MarshalJSON() ([]byte, error) {
	m := fiber.Map{
		"name":              service.Name,
		"url":               service.Url,
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

func (service *Service) check(httpClient *http.Client, checkInterval time.Duration) {
	for {
		startTime := time.Now()
		resp, err := httpClient.Head(service.Url)
		endTime := time.Now()

		if err != nil {
			logger.Println(err)

			service.ok = false
			service.responseTimeMilli = 0
			service.responseTimeMicro = 0
			service.error = err
		} else {
			if resp.StatusCode != 200 {
				err := fmt.Errorf("%v returned invalid status code %v", service.Url, resp.StatusCode)
				logger.Println(err)

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
