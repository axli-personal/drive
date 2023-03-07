package main

import (
	"github.com/axli-personal/drive/internal/drive/ports"
	"github.com/axli-personal/drive/internal/drive/service"
	"github.com/caarlos0/env/v7"
	"github.com/gofiber/fiber/v2"
	"net"
	"net/http"
	"net/rpc"
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

		app := fiber.New()

		app.Post("/drive/create", httpServer.CreateDrive)
		app.Post("/files/create", httpServer.CreateFile)
		app.Post("/folders/create", httpServer.CreateFolder)
		app.Get("/drive", httpServer.GetDrive)
		app.Get("/files/:fileId", httpServer.GetFile)
		app.Get("/folders/:folderId", httpServer.GetFolder)
		app.Get("/path/:parent", httpServer.GetPath)
		app.Post("/files/remove", httpServer.RemoveFile)

		err := app.Listen(":8080")
		if err != nil {
			panic(err)
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		rpcServer := ports.NewRPCServer(svc)

		err := rpc.Register(&rpcServer)
		if err != nil {
			panic(err)
		}

		rpc.HandleHTTP()

		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}

		err = http.Serve(listener, nil)
		if err != nil {
			panic(err)
		}
	}()

	waitGroup.Wait()
}
