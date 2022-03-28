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
		app.Get(utils.BaseExhibitionSearch, h.Search)
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
		page, err := strconv.Atoi(url.Get("page"))
		if err != nil {
			page = 1
		}
		size, err := strconv.Atoi(url.Get("size"))
		if err != nil {
			size = 10
		}
		exhibitionPage := h.u.GetExhibitions(page, size)
		return ctx.JSON(exhibitionPage)
	}
	exhibitions := h.u.GetExhibitionsByMuseum(museumID)
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

func (h *ExhibitionHandler) Search(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	content := h.u.Search(url.Get("name"))
	return ctx.JSON(content)
}
