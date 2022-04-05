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
		app.Get(utils.GatewayApiExhibitionID, h.GetExhibition)
		app.Get(utils.GatewayApiMuseumID, h.GetMuseum)
		app.Get(utils.GatewayApiMuseums, h.GetMuseums)
		app.Get(utils.GatewayApiExhibitions, h.GetExhibitions)
		app.Get(utils.GatewayApiSearch, h.Search)
		app.Get(utils.GatewayApiPictures, h.GetPictures)
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

func (h *GatewayHandler) GetExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	exhibition, msg := h.u.GetExhibition(id)
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(exhibition)
}

func (h *GatewayHandler) GetMuseum(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	museum, msg := h.u.GetMuseum(id)
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) GetMuseums(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	content := h.u.GetMuseums(url.RawQuery)
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetExhibitions(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if url.Query().Has("museumID") {
		exhibitions := h.u.GetMuseumExhibitions(url.RawQuery)
		return ctx.JSON(exhibitions)
	}
	content := h.u.GetExhibitions(url.RawQuery)
	return ctx.JSON(content)
}

func (h *GatewayHandler) Search(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	name := url.Get("name")
	if id, err := strconv.Atoi(url.Get("id")); err == nil {
		result := h.u.SearchByID(name, url.Get("type"), id)
		if result == nil {
			ctx.SetStatus(http.StatusNotFound)
			return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
		}
		return ctx.JSON(result)
	}

	content := h.u.Search(name)
	if content == nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
	}
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetPictures(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	ids := url.Get("id")
	exhibition := url.Get("exhibitionID")
	var pictures []*domain.Picture
	if exhibition != "" {
		pictures = h.u.GetPicturesExh(exhibition)
	} else {
		pictures = h.u.GetPicturesFav(ids)
	}
	if pictures == nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
	}
	return ctx.JSON(pictures)
}
