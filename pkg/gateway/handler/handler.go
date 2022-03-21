package handler

import (
	"ar_exhibition/pkg/gateway/usecase"
	"ar_exhibition/pkg/utils"

	"github.com/aerogo/aero"
)

type GatewayHandler struct {
	u *usecase.GatewayUsecase
}

func GatewayHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.GatewayUsecase)
	if ok {
		return &GatewayHandler{u: instance}
	}
	return nil
}

func ConfigureGateway(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*GatewayHandler)
	if ok {
		app.Get(utils.GatewayApiMain, h.GetMain)
	}
	return app
}

func (h *GatewayHandler) GetMain(ctx aero.Context) error {
	content := h.u.GetMain()
	return ctx.JSON(content)
}
