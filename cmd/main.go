package main

import (
	"log"

	"github.com/dkshi/chopchop"
	"github.com/dkshi/chopchop/internal/handler"
	"github.com/dkshi/chopchop/internal/repository"
	"github.com/dkshi/chopchop/internal/service"
)

const (
	port = "8080"
)

func main() {
	newRepository := repository.NewRepository()
	newService := service.NewService(newRepository)
	newHandler := handler.NewHandler(newService)
	newRouter := newHandler.InitRoutes()

	server := new(chopchop.Server)

	log.Printf("Server is now running on port: %s", port)

	if err := server.Run(port, newRouter); err != nil {
		log.Fatalf("error running the server %s", err)
	}
}
