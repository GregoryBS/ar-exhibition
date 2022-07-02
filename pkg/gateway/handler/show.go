package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

func (h *GatewayHandler) ShowMuseum(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header(utils.AuthHeader))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	err = h.u.ShowMuseum(id, user.ID)
	if err != nil {
		log.Println("Unable to publish museum with id:", id)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return nil
}

func (h *GatewayHandler) ShowExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header(utils.AuthHeader))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	err = h.u.ShowExhibition(id, user.ID)
	if err != nil {
		log.Println("Unable to publish exhibition with id:", id)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return nil
}

func (h *GatewayHandler) ShowPicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header(utils.AuthHeader))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	err = h.u.ShowPicture(id, user.ID)
	if err != nil {
		log.Println("Unable to publish picture with id:", id)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return nil
}
