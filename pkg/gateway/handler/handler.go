package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/gateway/usecase"
	"ar_exhibition/pkg/utils"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aerogo/aero"
)

const (
	PicsDir = "./pictures/"
)

type GatewayHandler struct {
	u *usecase.GatewayUsecase
}

func GatewayHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.GatewayUsecase)
	if ok {
		return &GatewayHandler{u: instance}
	}
	return nil
}

func ConfigureGateway(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*GatewayHandler)
	if ok {
		app.Get(utils.GatewayApiMain, h.GetMain)
		app.Get(utils.GatewayApiPictureID, h.GetPicture)
		app.Get(utils.GatewayApiExhibitionID, h.GetExhibition)
		app.Get(utils.GatewayApiMuseumID, h.GetMuseum)
		app.Get(utils.GatewayApiMuseums, h.GetMuseums)
		app.Get(utils.GatewayApiExhibitions, h.GetExhibitions)
		app.Get(utils.GatewayApiSearch, h.Search)
		app.Get(utils.GatewayApiPictures, h.GetPictures)
		app.Post(utils.GatewayApiMuseums, h.CreateMuseum)
		app.Post(utils.GatewayApiMuseumID, h.UpdateMuseum)
		app.Post(utils.GatewayApiMuseumImage, h.UpdateMuseumImage)
		app.Post(utils.GatewayApiPictures, h.CreatePicture)
		app.Post(utils.GatewayApiPictureImage, h.UpdatePictureImage)
	}
	return app
}

func (h *GatewayHandler) GetMain(ctx aero.Context) error {
	content := h.u.GetMain()
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetPicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	picture, msg := h.u.GetPicture(id)
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(picture)
}

func (h *GatewayHandler) GetExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	exhibition, msg := h.u.GetExhibition(id)
	if msg != nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(msg)
	}
	return ctx.JSON(exhibition)
}

func (h *GatewayHandler) GetMuseum(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
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
	content := h.u.GetMuseums(url.RawQuery)
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetExhibitions(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if url.Query().Has("museumID") {
		exhibitions := h.u.GetMuseumExhibitions(url.RawQuery)
		return ctx.JSON(exhibitions)
	}
	content := h.u.GetExhibitions(url.RawQuery)
	return ctx.JSON(content)
}

func (h *GatewayHandler) Search(ctx aero.Context) error {
	url := ctx.Request().Internal().URL
	if url.Query().Has("id") {
		result := h.u.SearchByID(url.Query().Get("type"), url.RawQuery)
		if result == nil {
			ctx.SetStatus(http.StatusNotFound)
			return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
		}
		return ctx.JSON(result)
	}
	content := h.u.Search(url.Query().Get("type"), url.RawQuery)
	if content == nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
	}
	return ctx.JSON(content)
}

func (h *GatewayHandler) GetPictures(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	ids := url.Get("id")
	exhibition := url.Get("exhibitionID")
	var pictures []*domain.Picture
	if exhibition != "" {
		pictures = h.u.GetPicturesExh(exhibition)
	} else {
		pictures = h.u.GetPicturesFav(ids)
	}
	if pictures == nil {
		ctx.SetStatus(http.StatusNotFound)
		return ctx.JSON(domain.ErrorResponse{Message: "Not found"})
	}
	return ctx.JSON(pictures)
}

func (h *GatewayHandler) CreateMuseum(ctx aero.Context) error {
	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to create"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	museum, err := h.u.CreateMuseum(museum, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) UpdateMuseum(ctx aero.Context) error {
	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	museum := new(domain.Museum)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), museum); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid museum to update"})
	}
	museum, err := h.u.UpdateMuseum(museum, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) UpdateMuseumImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	image, sizes := uploadImage(ctx.Request().Internal(), PicsDir)
	if image == "" {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload museum image"})
	}
	result := h.u.UpdateMuseumImage(image, sizes, id, user.ID)
	if result != nil {
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func (h *GatewayHandler) CreatePicture(ctx aero.Context) error {
	pic := new(domain.Picture)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), pic); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to create"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	pic, err := h.u.CreatePicture(pic, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(pic)
}

func (h *GatewayHandler) UpdatePictureImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		resp := domain.ErrorResponse{Message: "id not a number"}
		return ctx.JSON(resp)
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	image, sizes := uploadImage(ctx.Request().Internal(), PicsDir)
	if image == "" {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload picture image"})
	}
	result := h.u.UpdatePictureImage(image, sizes, id, user.ID)
	if result != nil {
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func checkAuth(header string) *domain.User {
	req, _ := http.NewRequest(http.MethodGet, utils.UserService+utils.UserID, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	user := new(domain.User)
	utils.DecodeJSON(resp.Body, user)
	return user
}

func uploadImage(r *http.Request, path string) (string, *domain.ImageSize) {
	r.ParseMultipartForm(10 << 20)
	reader, handler, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return "", nil
	}
	size := &domain.ImageSize{Height: m.Bounds().Dy(), Width: m.Bounds().Dx()}

	filename := utils.RandString(32) + filepath.Ext(handler.Filename)
	file, err := createFile(path, filename)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", nil
	}
	return filename, size
}

func createFile(dir, name string) (*os.File, error) {
	_, err := os.ReadDir(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(dir + name)
	return file, err
}
