package main

import (
	"ar_exhibition/pkg/exhibition/handler"
	"ar_exhibition/pkg/exhibition/repository"
	"ar_exhibition/pkg/exhibition/usecase"
	"ar_exhibition/pkg/server"
)

func main() {
	server.Run(handler.ConfigureExhibition,
		handler.ExhibitionHandlers,
		usecase.ExhibitionUsecases,
		repository.ExhibitionRepo)
}
