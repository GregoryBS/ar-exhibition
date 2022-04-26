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
		app.Post(utils.BaseExhibitionApi, h.Create)
		app.Post(utils.ExhibitionImage, h.UpdateImage)
		app.Post(utils.ExhibitionID, h.Update)
		app.Post(utils.ExhibitionShowID, h.ShowID)
		app.Delete(utils.ExhibitionID, h.Delete)
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
	var exhibitions []*domain.Exhibition
	if user, err := strconv.Atoi(ctx.Request().Header(utils.UserHeader)); err == nil {
		exhibitions = h.u.GetExhibitionsByUser(user)
	} else {
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
		exhibitions = h.u.GetExhibitionsByMuseum(museumID, filter)
	}
	return ctx.JSON(exhibitions)
}

func (h *ExhibitionHandler) GetExhibitionID(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	id, _ := strconv.Atoi(ctx.Get("id"))
	exhibition, err := h.u.GetExhibitionID(id, user)
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

func (h *ExhibitionHandler) ShowID(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	id, _ := strconv.Atoi(ctx.Get("id"))
	err := h.u.ShowID(id, user)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}

func (h *ExhibitionHandler) Create(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	data := new(domain.MuseumExhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), data); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to create"})
	}
	museum, exhibition := data.Mus, data.Exh

	exhibition = h.u.Create(exhibition, museum, user)
	if exhibition == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to create"})
	}
	return ctx.JSON(exhibition)
}

func (h *ExhibitionHandler) UpdateImage(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}

	exhibition = h.u.UpdateImage(exhibition, user)
	if exhibition == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	return ctx.JSON(exhibition)
}

func (h *ExhibitionHandler) Update(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}

	exhibition = h.u.Update(exhibition, user)
	if exhibition == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	return ctx.JSON(exhibition)
}

func (h *ExhibitionHandler) Delete(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	id, _ := strconv.Atoi(ctx.Get("id"))
	err := h.u.Delete(id, user)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}
