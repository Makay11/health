package health

import (
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
	Name              string    `json:"name"`
	Url               string    `json:"url"`
	Ok                bool      `json:"ok"`
	ResponseTimeMilli int64     `json:"responseTimeMilli"`
	ResponseTimeMicro int64     `json:"responseTimeMicro"`
	LastSeen          time.Time `json:"lastSeen"`
	Error             string    `json:"error,omitempty"`
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
			service.Error = err.Error()
		} else {
			if resp.StatusCode != 200 {
				err := fmt.Errorf("%v returned invalid status code %v", service.Url, resp.StatusCode)
				logger.Println(err)

				service.Ok = false
				service.Error = err.Error()
			} else {
				// logger.Printf("%v returned valid status code %v", service.Url, resp.StatusCode)

				service.Ok = true
				service.Error = ""
			}

			service.ResponseTimeMilli = endTime.UnixMilli() - startTime.UnixMilli()
			service.ResponseTimeMicro = endTime.UnixMicro() - startTime.UnixMicro()
			service.LastSeen = endTime
		}

		time.Sleep(checkInterval)
	}
}
