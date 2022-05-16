package main

import (
	"ar_exhibition/pkg/museum/handler"
	"ar_exhibition/pkg/museum/repository"
	"ar_exhibition/pkg/museum/usecase"
	"ar_exhibition/pkg/server"
)

func main() {
	server.Run(handler.ConfigureMuseum,
		nil,
		handler.MuseumHandlers,
		usecase.MuseumUsecases,
		repository.MuseumRepo)
}
