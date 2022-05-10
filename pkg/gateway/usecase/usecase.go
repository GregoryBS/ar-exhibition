package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type GatewayUsecase struct {
}

func GatewayUsecases(interface{}) interface{} {
	return &GatewayUsecase{}
}

func (u *GatewayUsecase) GetMain() *domain.MainPage {
	museums := make([]*domain.Museum, 0)
	resp, err := http.Get(utils.MuseumService + utils.MuseumTop)
	if err == nil {
		utils.DecodeJSON(resp.Body, &museums)
		resp.Body.Close()
	} else {
		log.Println("Cannot get museums for main:", err)
	}

	exhibitions := make([]*domain.Exhibition, 0)
	resp, err = http.Get(utils.ExhibitionService + utils.ExhibitionTop)
	if err == nil {
		utils.DecodeJSON(resp.Body, &exhibitions)
		resp.Body.Close()
	} else {
		log.Println("Cannot get exhibitions for main:", err)
	}

	pictures := make([]*domain.Picture, 0)
	resp, err = http.Get(utils.PictureService + utils.PictureTop)
	if err != nil {
		log.Println("Cannot get pictures for main:", err)
		return &domain.MainPage{Museums: museums, Exhibitions: exhibitions}
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, &pictures)
	return &domain.MainPage{Museums: museums, Exhibitions: exhibitions, Pictures: pictures}
}

func (u *GatewayUsecase) GetPicture(id int, user *domain.User) (*domain.Picture, *domain.ErrorResponse) {
	req, _ := http.NewRequest(http.MethodGet, utils.PictureService+strings.Replace(utils.PictureID, ":id", fmt.Sprint(id), 1), nil)
	if user != nil {
		req.Header.Set(utils.UserHeader, fmt.Sprint(user.ID))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get picture:", id, err)
		return nil, &domain.ErrorResponse{Message: err.Error()}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{Message: "Not found"}
	}

	picture := &domain.Picture{}
	utils.DecodeJSON(resp.Body, picture)
	return picture, nil
}

func (u *GatewayUsecase) GetExhibition(id int, user *domain.User) (*domain.Exhibition, *domain.ErrorResponse) {
	req, _ := http.NewRequest(http.MethodGet, utils.ExhibitionService+strings.Replace(utils.ExhibitionID, ":id", fmt.Sprint(id), 1), nil)
	if user != nil {
		req.Header.Set(utils.UserHeader, fmt.Sprint(user.ID))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get exhibition:", id, err)
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &domain.ErrorResponse{Message: "Not found"}
	}
	exhibition := &domain.Exhibition{}
	utils.DecodeJSON(resp.Body, exhibition)
	resp.Body.Close()

	exhibition.Pictures = make([]*domain.Picture, 0)
	req, _ = http.NewRequest(http.MethodGet, utils.PictureService+utils.PictureByExhibition+fmt.Sprint(exhibition.ID), nil)
	if user != nil {
		req.Header.Set(utils.UserHeader, fmt.Sprint(user.ID))
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get pictures for exhibition:", id, err)
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode == http.StatusOK {
		utils.DecodeJSON(resp.Body, &exhibition.Pictures)
	}
	defer resp.Body.Close()
	return exhibition, nil
}

func (u *GatewayUsecase) GetExhibitionPictures(exhibition int) []*domain.Picture {
	pictures := make([]*domain.Picture, 0)
	resp, err := http.Get(utils.PictureService + utils.PictureByExhibition + fmt.Sprint(exhibition))
	if err != nil {
		log.Println("Cannot get pictures for exhibition:", exhibition, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		utils.DecodeJSON(resp.Body, &pictures)
	}
	return pictures
}

func (u *GatewayUsecase) GetMuseum(id int) (*domain.Museum, *domain.ErrorResponse) {
	museum := &domain.Museum{}
	resp, err := http.Get(utils.MuseumService + strings.Replace(utils.MuseumID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		log.Println("Cannot get museum:", id, err)
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{Message: "Not found"}
	}
	utils.DecodeJSON(resp.Body, museum)
	resp.Body.Close()

	museum.Exhibitions = u.GetMuseumExhibitions("?museumID=" + fmt.Sprint(museum.ID))
	return museum, nil
}

func (u *GatewayUsecase) GetUserMuseums(user int) []*domain.Museum {
	req, _ := http.NewRequest(http.MethodGet, utils.MuseumService+utils.BaseMuseumApi, nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get user museums:", user, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	result := make([]*domain.Museum, 0)
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) GetMuseums(params string) *domain.Page {
	resp, err := http.Get(utils.MuseumService + utils.BaseMuseumApi + params)
	if err != nil {
		log.Println("Cannot get museums:", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	result := &domain.Page{}
	utils.DecodeJSON(resp.Body, result)
	return result
}

func (u *GatewayUsecase) GetMuseumExhibitions(params string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	resp, err := http.Get(utils.ExhibitionService + utils.BaseExhibitionApi + params)
	if err != nil {
		log.Println("Cannot get museum exhibitions:", err)
		return result
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) GetExhibitions(params string) *domain.Page {
	result := &domain.Page{}
	resp, err := http.Get(utils.ExhibitionService + utils.BaseExhibitionApi + params)
	if err != nil {
		log.Println("Cannot get exhibitions:", err)
		return result
	}
	defer resp.Body.Close()

	utils.DecodeJSON(resp.Body, result)
	return result
}

func (u *GatewayUsecase) GetUserExhibitions(user int) []*domain.Exhibition {
	req, _ := http.NewRequest(http.MethodGet, utils.ExhibitionService+utils.BaseExhibitionApi, nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get user exhibitions:", user, err)
		return nil
	}
	defer resp.Body.Close()

	result := make([]*domain.Exhibition, 0)
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) Search(param, params string) *domain.SearchPage {
	museums := make([]*domain.Museum, 0)
	exhibitions := make([]*domain.Exhibition, 0)
	pictures := make([]*domain.Picture, 0)
	switch param {
	case "museum":
		resp, err := http.Get(utils.MuseumService + utils.BaseMuseumSearch + params)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		utils.DecodeJSON(resp.Body, &museums)
		return &domain.SearchPage{Museums: museums}
	case "exhibition":
		resp, err := http.Get(utils.ExhibitionService + utils.BaseExhibitionSearch + params)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		utils.DecodeJSON(resp.Body, &exhibitions)
		return &domain.SearchPage{Exhibitions: exhibitions}
	case "picture":
		resp, err := http.Get(utils.PictureService + utils.BasePictureSearch + params)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		utils.DecodeJSON(resp.Body, &pictures)
		return &domain.SearchPage{Pictures: pictures}
	}
	resp, err := http.Get(utils.MuseumService + utils.BaseMuseumSearch + params)
	if err == nil {
		utils.DecodeJSON(resp.Body, &museums)
		resp.Body.Close()
	}

	resp, err = http.Get(utils.ExhibitionService + utils.BaseExhibitionSearch + params)
	if err == nil {
		utils.DecodeJSON(resp.Body, &exhibitions)
		resp.Body.Close()
	}

	resp, err = http.Get(utils.PictureService + utils.BasePictureSearch + params)
	if err == nil {
		utils.DecodeJSON(resp.Body, &pictures)
		resp.Body.Close()
	}
	return &domain.SearchPage{Museums: museums, Exhibitions: exhibitions, Pictures: pictures}
}

func (u *GatewayUsecase) SearchByID(param, params string) *domain.SearchPage {
	var resp *http.Response
	var err error
	switch param {
	case "exhibition":
		resp, err = http.Get(utils.ExhibitionService + utils.BaseExhibitionSearch + params)
	case "picture":
		resp, err = http.Get(utils.PictureService + utils.BasePictureSearch + params)
	default:
		return nil
	}
	if err != nil {
		log.Println("Search param:", param, "error:", err)
		return nil
	}
	defer resp.Body.Close()

	switch param {
	case "exhibition":
		result := make([]*domain.Exhibition, 0)
		utils.DecodeJSON(resp.Body, &result)
		return &domain.SearchPage{Exhibitions: result}
	case "picture":
		result := make([]*domain.Picture, 0)
		utils.DecodeJSON(resp.Body, &result)
		return &domain.SearchPage{Pictures: result}
	default:
		return nil
	}
}

func (u *GatewayUsecase) GetPicturesFav(id string) []*domain.Picture {
	resp, err := http.Get(utils.PictureService + utils.PictureByIDs + id)
	if err != nil {
		log.Println("Cannot get pictures for favourites:", id, err)
		return nil
	}
	defer resp.Body.Close()

	result := make([]*domain.Picture, 0)
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) GetPicturesUser(user int) []*domain.Picture {
	req, _ := http.NewRequest(http.MethodGet, utils.PictureService+utils.BasePictureApi, nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get pictures for user:", user, err)
		return nil
	}
	defer resp.Body.Close()

	result := make([]*domain.Picture, 0)
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) CreateMuseum(museum *domain.Museum, user int) (*domain.Museum, error) {
	req, _ := http.NewRequest(http.MethodPost, utils.MuseumService+utils.BaseMuseumApi,
		bytes.NewBuffer(utils.EncodeJSON(museum)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum creating error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to create museum")
	}
	result := new(domain.Museum)
	utils.DecodeJSON(resp.Body, result)
	return result, nil
}

func (u *GatewayUsecase) UpdateMuseum(museum *domain.Museum, user int) (*domain.Museum, error) {
	req, _ := http.NewRequest(http.MethodPost, utils.MuseumService+strings.Replace(utils.MuseumID, ":id", fmt.Sprint(museum.ID), 1),
		bytes.NewBuffer(utils.EncodeJSON(museum)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum updating error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to update museum")
	}
	result := new(domain.Museum)
	utils.DecodeJSON(resp.Body, result)
	return result, nil
}

func (u *GatewayUsecase) UpdateMuseumImage(filename string, sizes *domain.ImageSize, museum, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.MuseumService+strings.Replace(utils.MuseumImage, ":id", fmt.Sprint(museum), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Museum{ID: museum, Image: filename, Sizes: sizes})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum image updating error:", err)
		return &domain.ErrorResponse{Message: "Museum service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update museum image"}
	}
	return nil
}

func (u *GatewayUsecase) CreatePicture(pic *domain.Picture, user int) (*domain.Picture, error) {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+utils.BasePictureApi,
		bytes.NewBuffer(utils.EncodeJSON(pic)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture creating error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to create picture")
	}
	result := new(domain.Picture)
	utils.DecodeJSON(resp.Body, result)
	return result, nil
}

func (u *GatewayUsecase) UpdatePictureImage(filename string, sizes *domain.ImageSize, pic, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureImage, ":id", fmt.Sprint(pic), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Picture{ID: pic, Image: filename, Sizes: sizes})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture image updating error:", err)
		return &domain.ErrorResponse{Message: "Picture service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update picture image"}
	}
	return nil
}

func (u *GatewayUsecase) UpdatePictureVideo(filename, size string, pic, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureVideo, ":id", fmt.Sprint(pic), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Picture{ID: pic, Video: filename, VideoSize: size})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture video updating error:", err)
		return &domain.ErrorResponse{Message: "Picture service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update picture video"}
	}
	return nil
}

func (u *GatewayUsecase) UpdatePicture(picture *domain.Picture, user int) (*domain.Picture, error) {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureID, ":id", fmt.Sprint(picture.ID), 1),
		bytes.NewBuffer(utils.EncodeJSON(picture)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture updating error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to update picture")
	}
	result := new(domain.Picture)
	utils.DecodeJSON(resp.Body, result)
	return result, nil
}

func (u *GatewayUsecase) ShowMuseum(museum, user int) error {
	req, _ := http.NewRequest(http.MethodPost, utils.MuseumService+strings.Replace(utils.MuseumShow, ":id", fmt.Sprint(museum), 1), nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum publishing error:", museum, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to publish museum")
	}

	req, _ = http.NewRequest(http.MethodPost, utils.ExhibitionService+utils.ExhibitionShow, nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum exhibitions publishing error:", museum, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to publish museum exhibitions")
	}

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PictureShow, nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum pictures publishing error:", museum, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to publish museum pictures")
	}
	return nil
}

func (u *GatewayUsecase) ShowExhibition(exhibition, user int) error {
	req, _ := http.NewRequest(http.MethodPost, utils.ExhibitionService+strings.Replace(utils.ExhibitionShowID, ":id", fmt.Sprint(exhibition), 1), nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition publishing error:", exhibition, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to publish exhibition")
	}

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PictureShowExh+fmt.Sprint(exhibition), nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition pictures publishing error:", exhibition, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to publish exhibition pictures")
	}
	return nil
}

func (u *GatewayUsecase) ShowPicture(picture, user int) error {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureShowID, ":id", fmt.Sprint(picture), 1), nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture publishing error:", picture, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to publish picture")
	}
	return nil
}

func (u *GatewayUsecase) CreateExhibition(exhibition *domain.Exhibition, user int) (*domain.Exhibition, error) {
	museum := u.GetUserMuseums(user)
	if museum == nil {
		log.Println("Cannot get user museum for exhibition change")
		return nil, errors.New("Unable to create exhibition")
	}

	req, _ := http.NewRequest(http.MethodPost, utils.ExhibitionService+utils.BaseExhibitionApi,
		bytes.NewBuffer(utils.EncodeJSON(domain.MuseumExhibition{Mus: museum[0], Exh: exhibition})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition creating error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to create exhibition")
	}
	result := new(domain.Exhibition)
	utils.DecodeJSON(resp.Body, result)
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PicturesToExh,
		bytes.NewBuffer(utils.EncodeJSON(domain.MuseumExhibition{Mus: museum[0], Exh: result})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Pictures to exhibition adding error:", result.ID, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to create exhibition")
	}
	return result, nil
}

func (u *GatewayUsecase) UpdateExhibitionImage(filename string, sizes *domain.ImageSize, id, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.ExhibitionService+strings.Replace(utils.ExhibitionImage, ":id", fmt.Sprint(id), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Exhibition{ID: id, Image: filename, Sizes: sizes})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition image updating error:", id, err)
		return &domain.ErrorResponse{Message: "Exhibition service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update exhibition image"}
	}
	return nil
}

func (u *GatewayUsecase) UpdateExhibition(exhibition *domain.Exhibition, user int) (*domain.Exhibition, error) {
	museum := u.GetUserMuseums(user)
	if museum == nil {
		log.Println("Cannot get user museum for exhibition change")
		return nil, errors.New("Unable to update exhibition")
	}

	req, _ := http.NewRequest(http.MethodPost, utils.ExhibitionService+strings.Replace(utils.ExhibitionID, ":id", fmt.Sprint(exhibition.ID), 1),
		bytes.NewBuffer(utils.EncodeJSON(exhibition)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition updating error:", exhibition.ID, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to update exhibition")
	}
	result := new(domain.Exhibition)
	utils.DecodeJSON(resp.Body, result)
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PicturesToExh,
		bytes.NewBuffer(utils.EncodeJSON(domain.MuseumExhibition{Mus: museum[0], Exh: result})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Pictures to exhibition adding error:", result.ID, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to update exhibition")
	}
	return result, nil
}

func (u *GatewayUsecase) DeletePicture(id, user int) error {
	req, _ := http.NewRequest(http.MethodDelete, utils.PictureService+strings.Replace(utils.PictureID, ":id", fmt.Sprint(id), 1), nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture deleting error:", id, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to delete picture")
	}
	return nil
}

func (u *GatewayUsecase) DeleteExhibition(id, user int) error {
	req, _ := http.NewRequest(http.MethodDelete, utils.ExhibitionService+strings.Replace(utils.ExhibitionID, ":id", fmt.Sprint(id), 1), nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("exhibition deleting error:", id, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to delete exhibition")
	}

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PicturesToExh,
		bytes.NewBuffer(utils.EncodeJSON(domain.MuseumExhibition{Exh: &domain.Exhibition{ID: id}})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Pictures to exhibition deleting error:", id, err)
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to delete exhibition")
	}
	return nil
}
