package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/gateway/usecase"
	"ar_exhibition/pkg/utils"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aerogo/aero"
)

const (
	PicsDir  = "./pictures/"
	VideoDir = "./videos/"
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
		app.Post(utils.GatewayApiPictureVideo, h.UpdatePictureVideo)
		app.Post(utils.GatewayApiPictureID, h.UpdatePicture)
		app.Post(utils.GatewayApiMuseumShow, h.ShowMuseum)
		app.Post(utils.GatewayApiExhibitions, h.CreateExhibition)
		app.Post(utils.GatewayApiExhibitionID, h.UpdateExhibition)
		app.Post(utils.GatewayApiExhibitionImage, h.UpdateExhibitionImage)
		app.Post(utils.GatewayApiExhibitionShow, h.ShowExhibition)
		app.Post(utils.GatewayApiPictureShow, h.ShowPicture)
		app.Delete(utils.GatewayApiPictureID, h.DeletePicture)
		app.Delete(utils.GatewayApiExhibitionID, h.DeleteExhibition)
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	picture, msg := h.u.GetPicture(id, checkAuth(ctx.Request().Header("Authorization")))
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

	exhibition, msg := h.u.GetExhibition(id, checkAuth(ctx.Request().Header("Authorization")))
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
	if user := checkAuth(ctx.Request().Header("Authorization")); user != nil {
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
	if user := checkAuth(ctx.Request().Header("Authorization")); user != nil {
		if user.ID > 0 {
			exhibitions := h.u.GetUserExhibitions(user.ID)
			return ctx.JSON(exhibitions)
		}
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Authorization error"})
	} else if url.Query().Has("museumID") {
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
	var pictures []*domain.Picture
	if user := checkAuth(ctx.Request().Header("Authorization")); user != nil {
		if user.ID > 0 {
			pictures = h.u.GetPicturesUser(user.ID)
		}
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
	if user == nil || user.ID <= 0 {
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
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

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
	museum.ID = id
	museum, err = h.u.UpdateMuseum(museum, user.ID)
	if err != nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(museum)
}

func (h *GatewayHandler) UpdateMuseumImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload museum image"})
	}
	result := h.u.UpdateMuseumImage(image, sizes.(*domain.ImageSize), id, user.ID)
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
	if user == nil || user.ID <= 0 {
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload picture image"})
	}
	result := h.u.UpdatePictureImage(image, sizes.(*domain.ImageSize), id, user.ID)
	if result != nil {
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func (h *GatewayHandler) UpdatePictureVideo(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload picture video"})
	}
	result := h.u.UpdatePictureVideo(video, size.(string), id, user.ID)
	if result != nil {
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func checkAuth(header string) *domain.User {
	req, _ := http.NewRequest(http.MethodGet, utils.UserService+utils.UserID, nil)
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

func uploadFiles(r *http.Request) (string, interface{}) {
	mr, err := r.MultipartReader()
	if err != nil {
		return "", nil
	}
	var sizes interface{}
	files := make([]string, 0)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "image":
			filename := utils.RandString(32) + filepath.Ext(part.FileName())
			file, err := createFile(PicsDir, filename)
			if err != nil {
				fmt.Println("Error while creating file:", err.Error())
				continue
			}
			_, err = io.Copy(file, part)
			if err == nil {
				files = append(files, filename)
				if sizes == nil {
					file.Seek(0, 0)
					m, _, err := image.Decode(file)
					if err == nil {
						sizes = &domain.ImageSize{Height: m.Bounds().Dy(), Width: m.Bounds().Dx()}
					}
				}
				file.Close()
			}
		case "video":
			filename := utils.RandString(32) + filepath.Ext(part.FileName())
			file, err := createFile(VideoDir, filename)
			if err != nil {
				fmt.Println("Error while creating file:", err.Error())
				return "", nil
			}
			defer file.Close()
			_, err = io.Copy(file, part)
			if err != nil {
				return "", nil
			}
			files = append(files, filename)
		case "video_size":
			buf := make([]byte, 1024)
			s, err := part.Read(buf)
			if err == nil {
				sizes = string(buf[:s])
			}
		default:
			continue
		}
	}
	return strings.Join(files, ","), sizes
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

func (h *GatewayHandler) UpdatePicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid picture to update"})
	}
	picture.ID = id
	picture, err = h.u.UpdatePicture(picture, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(picture)
}

func (h *GatewayHandler) ShowMuseum(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	err = h.u.ShowMuseum(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return nil
}

func (h *GatewayHandler) ShowExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	err = h.u.ShowExhibition(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return nil
}

func (h *GatewayHandler) ShowPicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}

	err = h.u.ShowPicture(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return nil
}

func (h *GatewayHandler) CreateExhibition(ctx aero.Context) error {
	exhibition := new(domain.Exhibition)
	if err := utils.DecodeJSON(ctx.Request().Body().Reader(), exhibition); err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to create"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	exhibition, err := h.u.CreateExhibition(exhibition, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(exhibition)
}

func (h *GatewayHandler) UpdateExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Invalid exhibition to update"})
	}
	exhibition.ID = id
	exhibition, err = h.u.UpdateExhibition(exhibition, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(exhibition)
}

func (h *GatewayHandler) UpdateExhibitionImage(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
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
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "Unable to upload exhibition image"})
	}
	result := h.u.UpdateExhibitionImage(image, sizes.(*domain.ImageSize), id, user.ID)
	if result != nil {
		ctx.SetStatus(http.StatusBadRequest)
	}
	return ctx.JSON(result)
}

func (h *GatewayHandler) DeletePicture(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	err = h.u.DeletePicture(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	return nil
}

func (h *GatewayHandler) DeleteExhibition(ctx aero.Context) error {
	id, err := strconv.Atoi(ctx.Get("id"))
	if err != nil {
		ctx.SetStatus(http.StatusBadRequest)
		return ctx.JSON(domain.ErrorResponse{Message: "id not a number"})
	}

	user := checkAuth(ctx.Request().Header("Authorization"))
	if user == nil || user.ID <= 0 {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	err = h.u.DeleteExhibition(id, user.ID)
	if err != nil {
		ctx.SetStatus(http.StatusForbidden)
		return ctx.JSON(domain.ErrorResponse{Message: "Not Authorized"})
	}
	return nil
}
