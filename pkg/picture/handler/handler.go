package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/picture/usecase"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	log.Println("Unknown object instead of picture handler")
	return nil
}

func ConfigurePicture(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*PictureHandler)
	if ok {
		app.Get(utils.BasePictureApi, h.GetPictures)
		app.Get(utils.PictureTop, h.GetPictureTop)
		app.Get(utils.PictureID, h.GetPictureID)
		app.Get(utils.BasePictureSearch, h.Search)
		app.Post(utils.BasePictureApi, h.Create)
		app.Post(utils.PictureImage, h.UpdateImage)
		app.Post(utils.PictureVideo, h.UpdateVideo)
		app.Post(utils.PictureID, h.Update)
		app.Post(utils.PictureShow, h.Show)
		app.Post(utils.PictureShowID, h.ShowID)
		app.Delete(utils.PictureID, h.Delete)
		app.Post(utils.PicturesToExh, h.UpdateForExhibition)
	}
	return app
}

func (h *PictureHandler) GetPictures(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	var pictures []*domain.Picture
	if url.Has("exhibitionID") {
		if exhibitionID, _ := strconv.Atoi(url.Get("exhibitionID")); exhibitionID > 0 {
			pictures = h.u.GetPicturesByExh(exhibitionID, user)
		}
	} else if url.Has("id") {
		ids := url.Get("id")
		id := make([]int, 0)
		for _, str := range strings.Split(ids, ",") {
			num, _ := strconv.Atoi(str)
			id = append(id, num)
		}
		pictures = h.u.GetPicturesByIDs(id)
	} else {
		pictures = h.u.GetPicturesByUser(user)
	}
	return ctx.JSON(pictures)
}

func (h *PictureHandler) GetPictureTop(ctx aero.Context) error {
	pictures := h.u.GetPictureTop()
	return ctx.JSON(pictures)
}

func (h *PictureHandler) GetPictureID(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	id, _ := strconv.Atoi(ctx.Get("id"))
	picture, err := h.u.GetPictureID(id, user)
	if err != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(nil)
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) Search(ctx aero.Context) error {
	var content []*domain.Picture
	url := ctx.Request().Internal().URL.Query()
	name := url.Get("name")
	if id, err := strconv.Atoi(url.Get("id")); err == nil {
		content = h.u.SearchID(name, id)
	} else {
		content = h.u.Search(name)
	}
	return ctx.JSON(content)
}
