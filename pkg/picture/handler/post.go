package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

func (h *PictureHandler) Create(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		log.Println("Invalid picture json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to create"})
	}

	picture = h.u.Create(picture, user)
	if picture == nil {
		log.Println("Picture creating error")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to create"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) Update(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		log.Println("Invalid picture json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}

	picture = h.u.Update(picture, user)
	if picture == nil {
		log.Println("Picture updating error")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) UpdateImage(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		log.Println("Invalid picture json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}

	picture = h.u.UpdateImage(picture, user)
	if picture == nil {
		log.Println("Picture image updating error")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) UpdateVideo(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		log.Println("Invalid picture json", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}

	picture = h.u.UpdateVideo(picture, user)
	if picture == nil {
		log.Println("Picture video updating error")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) UpdateForExhibition(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	data := new(domain.MuseumExhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), data); err != nil {
		log.Println("Invalid exhibition json for updating pictures", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	museum, exhibition := data.Mus, data.Exh

	err := h.u.UpdateForExhibition(exhibition, museum, user)
	if err != nil {
		log.Println("Error while adding pictures to exhibition", err)
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(nil)
}

func (h *PictureHandler) Show(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	url := ctx.Request().Internal().URL.Query()
	exhibitionID, err := strconv.Atoi(url.Get("exhibitionID"))
	if err == nil {
		err = h.u.ShowExh(exhibitionID, user)
	} else {
		err = h.u.Show(user)
	}
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}

func (h *PictureHandler) ShowID(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	id, _ := strconv.Atoi(ctx.Get("id"))
	err := h.u.ShowID(id, user)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}

func (h *PictureHandler) Delete(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	id, _ := strconv.Atoi(ctx.Get("id"))
	err := h.u.Delete(id, user)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
	}
	return ctx.JSON(nil)
}
