package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

func (h *GatewayHandler) UpdateMuseum(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		log.Println("Invalid museum json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}
	museum.ID = id
	museum, err = h.u.UpdateMuseum(museum, user.ID)
	if err != nil {
		log.Println("Cannot update museum for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) UpdateMuseumImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	image, sizes := uploadFiles(ctx.Request().Internal())
	if image == "" {
		log.Println("image url is blank")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload museum image"})
	}
	result := h.u.UpdateMuseumImage(image, sizes.(*domain.ImageSize), id, user.ID)
	if result != nil {
		log.Println("Cannot update museum image for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func (h *GatewayHandler) UpdatePicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		log.Println("Invalid picture json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	picture.ID = id
	picture, err = h.u.UpdatePicture(picture, user.ID)
	if err != nil {
		log.Println("Cannot update picture with id:", id, "for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(picture)
}

func (h *GatewayHandler) UpdatePictureImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	image, sizes := uploadFiles(ctx.Request().Internal())
	if image == "" {
		log.Println("image url is blank")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload picture image"})
	}
	result := h.u.UpdatePictureImage(image, sizes.(*domain.ImageSize), id, user.ID)
	if result != nil {
		log.Println("Cannot update picture image with id:", id, "for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func (h *GatewayHandler) UpdatePictureVideo(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	video, size := uploadFiles(ctx.Request().Internal())
	if video == "" {
		log.Println("video url is blank")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload picture video"})
	}
	result := h.u.UpdatePictureVideo(video, size.(string), id, user.ID)
	if result != nil {
		log.Println("Cannot update picture video with id:", id, "for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func (h *GatewayHandler) UpdateExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		log.Println("Invalid exhibition json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	exhibition.ID = id
	exhibition, err = h.u.UpdateExhibition(exhibition, user.ID)
	if err != nil {
		log.Println("Cannot update exhibition with id:", id, "for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(exhibition)
}

func (h *GatewayHandler) UpdateExhibitionImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		log.Println("Bad url")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	image, sizes := uploadFiles(ctx.Request().Internal())
	if image == "" {
		log.Println("image url is blank")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload exhibition image"})
	}
	result := h.u.UpdateExhibitionImage(image, sizes.(*domain.ImageSize), id, user.ID)
	if result != nil {
		log.Println("Cannot update exhibition image with id:", id, "for user: ", user.ID)
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}
