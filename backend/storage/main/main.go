package main

import (
	"github.com/axli-personal/drive/backend/storage/ports"
	"github.com/axli-personal/drive/backend/storage/service"
	"github.com/caarlos0/env/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"sync"
)

func main() {
	config := service.Config{}

	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	svc, err := service.NewService(config)
	if err != nil {
		panic(err)
	}

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		httpServer := ports.NewHTTPServer(svc)

		app := fiber.New(
			fiber.Config{
				BodyLimit: 100 * 1024 * 1024,
			},
		)

		app.Use(cors.New(cors.Config{
			AllowCredentials: true,
		}))

		app.Post("/upload", httpServer.Upload)
		app.Get("/download/:fileId", httpServer.Download)

		err := app.Listen(":8080")
		if err != nil {
			panic(err)
		}
	}()

	waitGroup.Wait()
}
