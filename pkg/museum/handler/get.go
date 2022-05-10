package handler

import (
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/aerogo/aero"
)

func (h *MuseumHandler) GetMuseumTop(ctx aero.Context) error {
	museums := h.u.GetMuseumTop()
	return ctx.JSON(museums)
}

func (h *MuseumHandler) GetMuseumID(ctx aero.Context) error {
	id, _ := strconv.Atoi(ctx.Get("id"))
	museum, err := h.u.GetMuseumID(id)
	if err != nil {
		log.Println("Error while getting museum:", id, err)
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(nil)
	}
	return ctx.JSON(museum)
}

func (h *MuseumHandler) GetMuseums(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	page, err := strconv.Atoi(url.Get("page"))
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(url.Get("size"))
	if err != nil {
		size = 10
	}
	if user, err := strconv.Atoi(ctx.Request().Header(utils.UserHeader)); err == nil {
		items := h.u.GetUserMuseums(user)
		if items == nil {
			ctx.SetStatus(http.StatusForbidden)
		}
		return ctx.JSON(items)
	} else {
		content := h.u.GetMuseums(page, size)
		return ctx.JSON(content)
	}
}

func (h *MuseumHandler) Search(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	content := h.u.Search(url.Get("name"))
	return ctx.JSON(content)
}
