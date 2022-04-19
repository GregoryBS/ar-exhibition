package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"bytes"
	"errors"
	"fmt"
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
	if err != nil {
		return nil
	}
	utils.DecodeJSON(resp.Body, &museums)
	resp.Body.Close()

	exhibitions := make([]*domain.Exhibition, 0)
	resp, err = http.Get(utils.ExhibitionService + utils.ExhibitionTop)
	if err != nil {
		return &domain.MainPage{Museums: museums}
	}
	utils.DecodeJSON(resp.Body, &exhibitions)
	resp.Body.Close()

	pictures := make([]*domain.Picture, 0)
	resp, err = http.Get(utils.PictureService + utils.BasePictureApi)
	if err != nil {
		return &domain.MainPage{Museums: museums, Exhibitions: exhibitions}
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, &pictures)
	return &domain.MainPage{Museums: museums, Exhibitions: exhibitions, Pictures: pictures}
}

func (u *GatewayUsecase) GetPicture(id int) (*domain.Picture, *domain.ErrorResponse) {
	picture := &domain.Picture{}
	resp, err := http.Get(utils.PictureService + strings.Replace(utils.PictureID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{Message: "Not found"}
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, picture)
	return picture, nil
}

func (u *GatewayUsecase) GetExhibition(id int) (*domain.Exhibition, *domain.ErrorResponse) {
	exhibition := &domain.Exhibition{}
	resp, err := http.Get(utils.ExhibitionService + strings.Replace(utils.ExhibitionID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{Message: "Not found"}
	}
	utils.DecodeJSON(resp.Body, exhibition)
	resp.Body.Close()

	exhibition.Pictures = make([]*domain.Picture, 0)
	resp, err = http.Get(utils.PictureService + utils.PictureByExhibition + fmt.Sprint(exhibition.ID))
	if err != nil {
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode == http.StatusOK {
		utils.DecodeJSON(resp.Body, &exhibition.Pictures)
	}
	defer resp.Body.Close()
	return exhibition, nil
}

func (u *GatewayUsecase) GetMuseum(id int) (*domain.Museum, *domain.ErrorResponse) {
	museum := &domain.Museum{}
	resp, err := http.Get(utils.MuseumService + strings.Replace(utils.MuseumID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		return nil, &domain.ErrorResponse{Message: err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{Message: "Not found"}
	}
	utils.DecodeJSON(resp.Body, museum)
	resp.Body.Close()

	museum.Exhibitions = u.GetMuseumExhibitions("museumID=" + fmt.Sprint(museum.ID))
	return museum, nil
}

func (u *GatewayUsecase) GetMuseums(params string) *domain.Page {
	result := &domain.Page{}
	resp, err := http.Get(utils.MuseumService + utils.BaseMuseumApi + "?" + params)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	utils.DecodeJSON(resp.Body, result)
	return result
}

func (u *GatewayUsecase) GetMuseumExhibitions(params string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	resp, err := http.Get(utils.ExhibitionService + utils.BaseExhibitionApi + "?" + params)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) GetExhibitions(params string) *domain.Page {
	result := &domain.Page{}
	resp, err := http.Get(utils.ExhibitionService + utils.BaseExhibitionApi + "?" + params)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	utils.DecodeJSON(resp.Body, result)
	return result
}

func (u *GatewayUsecase) Search(param, params string) *domain.SearchPage {
	museums := make([]*domain.Museum, 0)
	exhibitions := make([]*domain.Exhibition, 0)
	pictures := make([]*domain.Picture, 0)
	switch param {
	case "museum":
		resp, err := http.Get(utils.MuseumService + utils.BaseMuseumSearch + "?" + params)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		utils.DecodeJSON(resp.Body, &museums)
		return &domain.SearchPage{Museums: museums}
	case "exhibition":
		resp, err := http.Get(utils.ExhibitionService + utils.BaseExhibitionSearch + "?" + params)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		utils.DecodeJSON(resp.Body, &exhibitions)
		return &domain.SearchPage{Exhibitions: exhibitions}
	case "picture":
		resp, err := http.Get(utils.PictureService + utils.BasePictureSearch + "?" + params)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		utils.DecodeJSON(resp.Body, &pictures)
		return &domain.SearchPage{Pictures: pictures}
	}
	resp, err := http.Get(utils.MuseumService + utils.BaseMuseumSearch + "?" + params)
	if err != nil {
		return nil
	}
	utils.DecodeJSON(resp.Body, &museums)
	resp.Body.Close()

	resp, err = http.Get(utils.ExhibitionService + utils.BaseExhibitionSearch + "?" + params)
	if err != nil {
		return &domain.SearchPage{Museums: museums}
	}
	utils.DecodeJSON(resp.Body, &exhibitions)
	resp.Body.Close()

	resp, err = http.Get(utils.PictureService + utils.BasePictureSearch + "?" + params)
	if err != nil {
		return &domain.SearchPage{Museums: museums, Exhibitions: exhibitions}
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, &pictures)
	return &domain.SearchPage{Museums: museums, Exhibitions: exhibitions, Pictures: pictures}
}

func (u *GatewayUsecase) SearchByID(param, params string) *domain.SearchPage {
	var resp *http.Response
	var err error
	switch param {
	case "exhibition":
		resp, err = http.Get(utils.ExhibitionService + utils.BaseExhibitionSearch + "?" + params)
	case "picture":
		resp, err = http.Get(utils.PictureService + utils.BasePictureSearch + "?" + params)
	default:
		return nil
	}
	if err != nil {
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

func (u *GatewayUsecase) GetPicturesExh(id string) []*domain.Picture {
	resp, err := http.Get(utils.PictureService + utils.PictureByExhibition + id)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	result := make([]*domain.Picture, 0)
	utils.DecodeJSON(resp.Body, &result)
	return result
}

func (u *GatewayUsecase) GetPicturesFav(id string) []*domain.Picture {
	resp, err := http.Get(utils.PictureService + utils.PictureByIDs + id)
	if err != nil {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &domain.ErrorResponse{Message: "Picture service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update picture image"}
	}
	return nil
}
