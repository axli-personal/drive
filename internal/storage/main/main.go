package main

import (
	"github.com/axli-personal/drive/internal/storage/ports"
	"github.com/axli-personal/drive/internal/storage/service"
	"github.com/caarlos0/env/v7"
	"github.com/gofiber/fiber/v2"
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

		app.Post("/objects/:fileId", httpServer.UploadObject)
		app.Get("/objects/:fileId", httpServer.GetObject)

		err := app.Listen(":8080")
		if err != nil {
			panic(err)
		}
	}()

	waitGroup.Wait()
}
