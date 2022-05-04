package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/picture/usecase"
	"ar_exhibition/pkg/utils"
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

func (h *PictureHandler) Create(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to create"})
	}

	picture = h.u.Create(picture, user)
	if picture == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to create"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) UpdateImage(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}

	picture = h.u.UpdateImage(picture, user)
	if picture == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) UpdateVideo(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}

	picture = h.u.UpdateVideo(picture, user)
	if picture == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	return ctx.JSON(picture)
}

func (h *PictureHandler) Update(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	picture := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), picture); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}

	picture = h.u.Update(picture, user)
	if picture == nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	return ctx.JSON(picture)
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

func (h *PictureHandler) UpdateForExhibition(ctx aero.Context) error {
	user, _ := strconv.Atoi(ctx.Request().Header(utils.UserHeader))
	data := new(domain.MuseumExhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), data); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	museum, exhibition := data.Mus, data.Exh

	err := h.u.UpdateForExhibition(exhibition, museum, user)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(nil)
}
