package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/user"
	"ar_exhibition/pkg/user/usecase"
	"ar_exhibition/pkg/utils"
	"bytes"
	"log"
	"net/http"

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
	log.Println("Unknown object instead of user handler")
	return nil
}

func ConfigureUser(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*UserHandler)
	if ok {
		app.Post(utils.UserSignup, h.Signup)
		app.Post(utils.UserLogin, h.Login)
		app.Get(utils.UserID, h.Check)
	}
	return app
}

func (h *UserHandler) Signup(ctx aero.Context) error {
	form := &domain.User{}
	if utils.DecodeJSON(ctx.Request().Body().Reader(), form) != nil {
		log.Println("Invalid user json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid signup form"})
	}

	created, err := h.u.Signup(form)
	if err != nil {
		log.Println("Unable to create user:", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(err)
	}

	token, err := user.CreateJWT(created.ID)
	if err != nil {
		log.Println("Unexpected error of creating jwt")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Cannot create jwt-token"})
	}
	created.Token = token

	req, _ := http.NewRequest(http.MethodPost, utils.GatewayService+utils.GatewayApiMuseums,
		bytes.NewBuffer(utils.EncodeJSON(&domain.Museum{Name: form.Museum})))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error while creating museum for user:", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Cannot create museum for user"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Cannot create museum for user"})
	}
	return ctx.JSON(created)
}

func (h *UserHandler) Login(ctx aero.Context) error {
	form := &domain.User{}
	if utils.DecodeJSON(ctx.Request().Body().Reader(), form) != nil {
		log.Println("Invalid user json")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid login form"})
	}

	created, err := h.u.Login(form)
	if err != nil {
		log.Println("Unable to login user:", err)
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(err)
	}

	token, err := user.CreateJWT(created.ID)
	if err != nil {
		log.Println("Unexpected error of creating jwt")
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Cannot create jwt-token"})
	}
	created.Token = token
	return ctx.JSON(created)
}

func (h *UserHandler) Check(ctx aero.Context) error {
	id := user.CheckJWT(ctx.Request().Header("Authorization"))
	if id != 0 {
		return ctx.JSON(domain.User{ID: id})
	}
	ctx.SetStatus(http.StatusUnauthorized)
	return ctx.JSON(domain.ErrorResponse{Message: "user not authorized"})
}
