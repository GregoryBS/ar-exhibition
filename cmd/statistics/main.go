package main

import (
	"ar_exhibition/pkg/server"
	"ar_exhibition/pkg/statistics/handler"
	"ar_exhibition/pkg/statistics/repository"
	"ar_exhibition/pkg/statistics/usecase"
)

func main() {
	server.Run(handler.ConfigureStat,
		handler.ReceiveStats,
		handler.StatHandlers,
		usecase.StatUsecases,
		repository.StatRepo)
}
