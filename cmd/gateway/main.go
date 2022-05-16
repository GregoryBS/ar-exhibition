package main

import (
	"ar_exhibition/pkg/gateway/handler"
	"ar_exhibition/pkg/gateway/usecase"
	"ar_exhibition/pkg/server"
)

func main() {
	server.Run(handler.ConfigureGateway,
		nil,
		handler.GatewayHandlers,
		usecase.GatewayUsecases)
}
