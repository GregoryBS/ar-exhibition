package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/museum/usecase"
	"ar_exhibition/pkg/utils"
	"log"
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
	log.Println("Unknown object instead of museum handler")
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

func (h *MuseumHandler) Create(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		log.Println("Invalid museum json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to create"})
	}

	museum = h.u.Create(museum, user)
	if museum == nil {
		log.Println("Unable to create museum")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to create"})
	}
	return ctx.JSON(museum)
}

func (h *MuseumHandler) Update(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		log.Println("Invalid museum json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}

	museum = h.u.Update(museum, user)
	if museum == nil {
		log.Println("Unable to update museum")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}
	return ctx.JSON(museum)
}

func (h *MuseumHandler) UpdateImage(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		log.Println("Invalid museum json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}

	museum = h.u.UpdateImage(museum, user)
	if museum == nil {
		log.Println("Unable to update museum image")
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
		log.Println("Error while publishing museum:", err)
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}
