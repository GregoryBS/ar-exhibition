package handler

import (
	"ar_exhibition/pkg/exhibition/usecase"
	"ar_exhibition/pkg/utils"

	"github.com/aerogo/aero"
)

type ExhibitionHandler struct {
	u *usecase.ExhibitionUsecase
}

func ExhibitionHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.ExhibitionUsecase)
	if ok {
		return &ExhibitionHandler{u: instance}
	}
	return nil
}

func ConfigureExhibition(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*ExhibitionHandler)
	if ok {
		app.Get(utils.ExhibitionTop, h.GetExhibitionTop)
	}
	return app
}

func (h *ExhibitionHandler) GetExhibitionTop(ctx aero.Context) error {
	exhibitions := h.u.GetExhibitionTop()
	return ctx.JSON(exhibitions)
}
