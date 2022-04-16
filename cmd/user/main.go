package main

import (
	"ar_exhibition/pkg/server"
	"ar_exhibition/pkg/user/handler"
	"ar_exhibition/pkg/user/repository"
	"ar_exhibition/pkg/user/usecase"
)

func main() {
	server.Run(handler.ConfigureUser,
		handler.UserHandlers,
		usecase.UserUsecases,
		repository.UserRepo)
}
