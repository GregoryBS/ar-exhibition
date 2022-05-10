package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"

	"github.com/aerogo/aero"
)

func (h *GatewayHandler) CreateMuseum(ctx aero.Context) error {
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		log.Println("Invalid museum json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to create"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	museum, err := h.u.CreateMuseum(museum, user.ID)
	if err != nil {
		log.Println("Cannot create museum for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) CreatePicture(ctx aero.Context) error {
	pic := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), pic); err != nil {
		log.Println("Invalid picture json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to create"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	pic, err := h.u.CreatePicture(pic, user.ID)
	if err != nil {
		log.Println("Cannot create picture for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(pic)
}

func (h *GatewayHandler) CreateExhibition(ctx aero.Context) error {
	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		log.Println("Invalid exhibition json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to create"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	exhibition, err := h.u.CreateExhibition(exhibition, user.ID)
	if err != nil {
		log.Println("Cannot create exhibition for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(exhibition)
}
