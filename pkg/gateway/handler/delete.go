package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

func (h *GatewayHandler) DeletePicture(ctx aero.Context) error {
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
	err = h.u.DeletePicture(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	return nil
}

func (h *GatewayHandler) DeleteExhibition(ctx aero.Context) error {
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
	err = h.u.DeleteExhibition(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	return nil
}
