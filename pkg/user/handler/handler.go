package handler

import (
	"ar_exhibition/pkg/user/usecase"
	"ar_exhibition/pkg/utils"

	"github.com/aerogo/aero"
)

type UserHandler struct {
	u *usecase.UserUsecase
}

func UserHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.UserUsecase)
	if ok {
		return &UserHandler{u: instance}
	}
	return nil
}

func ConfigureUser(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*UserHandler)
	if ok {
		app.Post(utils.UserSignup, h.Signup)
		app.Post(utils.UserLogin, h.Login)
	}
	return app
}

func (h *UserHandler) Signup(ctx aero.Context) error {
	return nil
}

func (h *UserHandler) Login(ctx aero.Context) error {
	return nil
}
