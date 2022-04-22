package handler

import (
	"ar_exhibition/pkg/domain"
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
		app.Post(utils.ExhibitionShow, h.Show)
	}
	return app
}

func (h *ExhibitionHandler) GetExhibitionTop(ctx aero.Context) error {
	exhibitions := h.u.GetExhibitionTop()
	return ctx.JSON(exhibitions)
}

func (h *ExhibitionHandler) GetExhibitions(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	filter := url.Get("filter")
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
		exhibitionPage := h.u.GetExhibitions(page, size, filter)
		return ctx.JSON(exhibitionPage)
	}
	exhibitions := h.u.GetExhibitionsByMuseum(museumID, filter)
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
	var content []*domain.Exhibition
	url := ctx.Request().Internal().URL.Query()
	filter := url.Get("filter")
	name := url.Get("name")
	if id, err := strconv.Atoi(url.Get("id")); err == nil {
		content = h.u.SearchID(name, id, filter)
	} else {
		content = h.u.Search(name, filter)
	}
	return ctx.JSON(content)
}

func (h *ExhibitionHandler) Show(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	err := h.u.Show(user)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}
