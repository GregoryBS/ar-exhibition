package handler

import (
	"ar_exhibition/pkg/museum/usecase"
	"ar_exhibition/pkg/utils"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

type MuseumHandler struct {
	u *usecase.MuseumUsecase
}

func MuseumHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.MuseumUsecase)
	if ok {
		return &MuseumHandler{u: instance}
	}
	return nil
}

func ConfigureMuseum(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*MuseumHandler)
	if ok {
		app.Get(utils.MuseumTop, h.GetMuseumTop)
		app.Get(utils.MuseumID, h.GetMuseumID)
	}
	return app
}

func (h *MuseumHandler) GetMuseumTop(ctx aero.Context) error {
	museums := h.u.GetMuseumTop()
	return ctx.JSON(museums)
}

func (h *MuseumHandler) GetMuseumID(ctx aero.Context) error {
	id, _ := strconv.Atoi(ctx.Get("id"))
	museum, err := h.u.GetMuseumID(id)
	if err != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(nil)
	}
	return ctx.JSON(museum)
}
