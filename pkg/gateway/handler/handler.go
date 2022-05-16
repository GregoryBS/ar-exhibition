package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/gateway/usecase"
	"ar_exhibition/pkg/utils"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	log.Println("Unknown object instead of gateway handler")
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
		app.Get(utils.GatewayApiStats, h.GetStats)
	}
	return app
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
				log.Println("Error while creating file for image:", err.Error())
				continue
			}
			_, err = io.Copy(file, part)
			if err == nil {
				files = append(files, filename)
				if sizes == nil {
					file.Seek(0, 0)
					m, _, err := image.Decode(file)
					if err == nil {
						sizes = &domain.ImageSize{Height: float32(m.Bounds().Dy()), Width: float32(m.Bounds().Dx())}
					}
				}
				file.Close()
			}
		case "image_url":
			buf := make([]byte, 1024)
			if s, err := part.Read(buf); err == nil {
				urls := strings.Split(string(buf[:s]), ",")
				for i := range urls {
					urls[i] = urls[i][strings.LastIndex(urls[i], "/")+1:]
				}
				files = append(files, urls...)
			}
		case "video":
			filename := utils.RandString(32) + filepath.Ext(part.FileName())
			file, err := createFile(VideoDir, filename)
			if err != nil {
				log.Println("Error while creating file for video:", err.Error())
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
			if s, err := part.Read(buf); err == nil {
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
