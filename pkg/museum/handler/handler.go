package handler

import (
	"ar_exhibition/pkg/domain"
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
		app.Get(utils.BaseMuseumApi, h.GetMuseums)
		app.Get(utils.BaseMuseumSearch, h.Search)
		app.Post(utils.BaseMuseumApi, h.Create)
		app.Post(utils.MuseumID, h.Update)
		app.Post(utils.MuseumImage, h.UpdateImage)
		app.Post(utils.MuseumShow, h.Show)
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

func (h *MuseumHandler) GetMuseums(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	page, err := strconv.Atoi(url.Get("page"))
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(url.Get("size"))
	if err != nil {
		size = 10
	}
	var content *domain.Page
	if user, err := strconv.Atoi(ctx.Request().Header(utils.UserHeader)); err == nil {
		items := h.u.GetUserMuseums(user)
		if items == nil {
			ctx.SetStatus(http.StatusForbidden)
			return ctx.JSON(nil)
		}
		content = &domain.Page{Items: items}
	} else {
		content = h.u.GetMuseums(page, size)
	}
	return ctx.JSON(content)
}

func (h *MuseumHandler) Search(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	content := h.u.Search(url.Get("name"))
	return ctx.JSON(content)
}

func (h *MuseumHandler) Create(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to create"})
	}

	museum = h.u.Create(museum, user)
	if museum == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to create"})
	}
	return ctx.JSON(museum)
}

func (h *MuseumHandler) Update(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}

	museum = h.u.Update(museum, user)
	if museum == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}
	return ctx.JSON(museum)
}

func (h *MuseumHandler) UpdateImage(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}

	museum = h.u.UpdateImage(museum, user)
	if museum == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}
	return ctx.JSON(museum)
}

func (h *MuseumHandler) Show(ctx aero.Context) error {
	id, _ := strconv.Atoi(ctx.Get("id"))
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	err := h.u.Show(id, user)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}
