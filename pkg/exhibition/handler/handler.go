package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/exhibition/usecase"
	"ar_exhibition/pkg/utils"
	"log"
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
	log.Println("Unknown object instead of exhibition handler")
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
		log.Println("Invalid exhibition json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to create"})
	}
	museum, exhibition := data.Mus, data.Exh

	exhibition = h.u.Create(exhibition, museum, user)
	if exhibition == nil {
		log.Println("Unable to create exhibition")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to create"})
	}
	return ctx.JSON(exhibition)
}

func (h *ExhibitionHandler) UpdateImage(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		log.Println("Invalid exhibition json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}

	exhibition = h.u.UpdateImage(exhibition, user)
	if exhibition == nil {
		log.Println("Unable to update exhibition image")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	return ctx.JSON(exhibition)
}

func (h *ExhibitionHandler) Update(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		log.Println("Invalid exhibition json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}

	exhibition = h.u.Update(exhibition, user)
	if exhibition == nil {
		log.Println("Unable to update exhibition")
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
