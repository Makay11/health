package main

import (
	"log"
	"time"

	"github.com/Makay11/health"
	"github.com/gofiber/fiber/v2"
)

func main() {
	go startTestApi()

	time.Sleep(time.Second)

	err := health.Start(
		":80",
		"/",
		30*time.Second,
		10*time.Second,
		[]health.Service{
			{
				Name: "status",
				Url:  "http://localhost",
			},
			{
				Name: "fast",
				Url:  "http://localhost:3000/api/fast",
			},
			{
				Name: "medium",
				Url:  "http://localhost:3000/api/medium",
			},
			{
				Name: "slow",
				Url:  "http://localhost:3000/api/slow",
			},
		},
	)

	log.Fatalln(err)
}

func startTestApi() {
	app := fiber.New()

	app.Get("/api/fast", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/api/medium", func(c *fiber.Ctx) error {
		time.Sleep(1 * time.Second)
		return c.SendStatus(200)
	})

	app.Get("/api/slow", func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second)
		return c.SendStatus(200)
	})

	err := app.Listen(":3000")

	log.Fatalln(err)
}
