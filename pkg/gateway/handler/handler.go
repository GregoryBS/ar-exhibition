package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/gateway/usecase"
	"ar_exhibition/pkg/utils"
	"net/http"
	"strconv"

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
		app.Get(utils.GatewayApiPictureID, h.GetPicture)
	}
	return app
}

func (h *GatewayHandler) GetMain(ctx aero.Context) error {
	content := h.u.GetMain()
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetPicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	picture, msg := h.u.GetPicture(id)
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(picture)
}
