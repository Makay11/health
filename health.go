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
	Ok                bool
	ResponseTimeMilli int64
	ResponseTimeMicro int64
	LastSeen          time.Time
	Error             error
}

func (service *Service) MarshalJSON() ([]byte, error) {
	m := fiber.Map{
		"name":              service.Name,
		"url":               service.Url,
		"ok":                service.Ok,
		"responseTimeMilli": service.ResponseTimeMilli,
		"responseTimeMicro": service.ResponseTimeMicro,
	}

	if service.LastSeen.IsZero() {
		m["lastSeen"] = nil
	} else {
		m["lastSeen"] = service.LastSeen
	}

	if service.Error == nil {
		m["error"] = nil
	} else {
		m["error"] = service.Error.Error()
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

			service.Ok = false
			service.ResponseTimeMilli = 0
			service.ResponseTimeMicro = 0
			service.Error = err
		} else {
			if resp.StatusCode != 200 {
				err := fmt.Errorf("%v returned invalid status code %v", service.Url, resp.StatusCode)
				logger.Println(err)

				service.Ok = false
				service.Error = err
			} else {
				// logger.Printf("%v returned valid status code %v", service.Url, resp.StatusCode)

				service.Ok = true
				service.Error = nil
			}

			service.ResponseTimeMilli = endTime.UnixMilli() - startTime.UnixMilli()
			service.ResponseTimeMicro = endTime.UnixMicro() - startTime.UnixMicro()
			service.LastSeen = endTime
		}

		time.Sleep(checkInterval)
	}
}
