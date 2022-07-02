package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

func (h *GatewayHandler) GetMain(ctx aero.Context) error {
	content := h.u.GetMain()
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetPicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	picture, msg := h.u.GetPicture(id, checkAuth(ctx.Request().Header(utils.AuthHeader)))
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(picture)
}

func (h *GatewayHandler) GetExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	exhibition, msg := h.u.GetExhibition(id, checkAuth(ctx.Request().Header(utils.AuthHeader)))
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(exhibition)
}

func (h *GatewayHandler) GetMuseum(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	museum, msg := h.u.GetMuseum(id)
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) GetMuseums(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if user := checkAuth(ctx.Request().Header(utils.AuthHeader)); user != nil {
		if user.ID > 0 {
			museums := h.u.GetUserMuseums(user.ID)
			if museums == nil {
				ctx.SetStatus(http.StatusNotFound)
				return ctx.JSON(domain.ErrorResponse{Message: "Cannot find museums for request"})
			}
			return ctx.JSON(museums[0])
		}
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Authorization error"})
	}
	content := h.u.GetMuseums("?" + url.RawQuery)
	if content != nil {
		return ctx.JSON(content)
	}
	ctx.SetStatus(http.StatusNotFound)
	return ctx.JSON(domain.ErrorResponse{Message: "Cannot find museums for request"})
}

func (h *GatewayHandler) GetExhibitions(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if user := checkAuth(ctx.Request().Header(utils.AuthHeader)); user != nil {
		if user.ID > 0 {
			exhibitions := h.u.GetUserExhibitions(user.ID)
			return ctx.JSON(exhibitions)
		}
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Authorization error"})
	} else if url.Query().Has("museumID") {
		exhibitions := h.u.GetMuseumExhibitions("?" + url.RawQuery)
		return ctx.JSON(exhibitions)
	}
	content := h.u.GetExhibitions("?" + url.RawQuery)
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetPictures(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	ids := url.Get("id")
	var pictures []*domain.Picture
	if user := checkAuth(ctx.Request().Header(utils.AuthHeader)); user != nil {
		if user.ID > 0 {
			pictures = h.u.GetPicturesUser(user.ID)
		}
	} else if url.Has("exhibitionID") {
		exhibition, _ := strconv.Atoi(url.Get("exhibitionID"))
		pictures = h.u.GetExhibitionPictures(exhibition)
	} else {
		pictures = h.u.GetPicturesFav(ids)
	}
	if pictures == nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
	}
	return ctx.JSON(pictures)
}

func (h *GatewayHandler) Search(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if url.Query().Has("id") {
		result := h.u.SearchByID(url.Query().Get("type"), "?"+url.RawQuery)
		if result == nil {
			ctx.SetStatus(http.StatusNotFound)
			return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
		}
		return ctx.JSON(result)
	}
	content := h.u.Search(url.Query().Get("type"), "?"+url.RawQuery)
	if content == nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
	}
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetStats(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if user := checkAuth(ctx.Request().Header(utils.AuthHeader)); user != nil {
		result := h.u.GetStats(user, "?"+url.RawQuery)
		return ctx.JSON(result)
	}
	ctx.SetStatus(http.StatusForbidden)
	return nil
}
