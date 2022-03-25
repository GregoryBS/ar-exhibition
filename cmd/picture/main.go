package main

import (
	"ar_exhibition/pkg/picture/handler"
	"ar_exhibition/pkg/picture/repository"
	"ar_exhibition/pkg/picture/usecase"
	"ar_exhibition/pkg/server"
)

func main() {
	server.Run(handler.ConfigurePicture,
		handler.PictureHandlers,
		usecase.PictureUsecases,
		repository.PictureRepo)
}
