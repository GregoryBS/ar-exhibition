package handler

import (
	"ar_exhibition/pkg/picture/usecase"
	"ar_exhibition/pkg/utils"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

type PictureHandler struct {
	u *usecase.PictureUsecase
}

func PictureHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.PictureUsecase)
	if ok {
		return &PictureHandler{u: instance}
	}
	return nil
}

func ConfigurePicture(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*PictureHandler)
	if ok {
		app.Get(utils.BasePictureApi, h.GetPictures)
		app.Get(utils.PictureID, h.GetPictureID)
	}
	return app
}

func (h *PictureHandler) GetPictures(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	exhibitionID, err := strconv.Atoi(url.Get("exhibitionID"))
	if err != nil {
		exhibitionID = 0
	}
	pictures := h.u.GetPictures(exhibitionID)
	return ctx.JSON(pictures)
}

func (h *PictureHandler) GetPictureID(ctx aero.Context) error {
	id, _ := strconv.Atoi(ctx.Get("id"))
	picture, err := h.u.GetPictureID(id)
	if err != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(nil)
	}
	return ctx.JSON(picture)
}
