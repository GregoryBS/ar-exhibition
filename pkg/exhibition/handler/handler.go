package handler

import (
	"ar_exhibition/pkg/exhibition/usecase"
	"ar_exhibition/pkg/utils"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

type ExhibitionHandler struct {
	u *usecase.ExhibitionUsecase
}

func ExhibitionHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.ExhibitionUsecase)
	if ok {
		return &ExhibitionHandler{u: instance}
	}
	return nil
}

func ConfigureExhibition(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*ExhibitionHandler)
	if ok {
		app.Get(utils.ExhibitionTop, h.GetExhibitionTop)
		app.Get(utils.BaseExhibitionApi, h.GetExhibitions)
		app.Get(utils.ExhibitionID, h.GetExhibitionID)
	}
	return app
}

func (h *ExhibitionHandler) GetExhibitionTop(ctx aero.Context) error {
	exhibitions := h.u.GetExhibitionTop()
	return ctx.JSON(exhibitions)
}

func (h *ExhibitionHandler) GetExhibitions(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	museumID, err := strconv.Atoi(url.Get("museumID"))
	if err != nil {
		museumID = 0
	}
	exhibitions := h.u.GetExhibitions(museumID)
	return ctx.JSON(exhibitions)
}

func (h *ExhibitionHandler) GetExhibitionID(ctx aero.Context) error {
	id, _ := strconv.Atoi(ctx.Get("id"))
	exhibition, err := h.u.GetExhibitionID(id)
	if err != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(nil)
	}
	return ctx.JSON(exhibition)
}
