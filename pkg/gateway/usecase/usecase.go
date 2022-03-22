package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"fmt"
	"net/http"
	"strings"
)

const (
	timeLayout = "2006-01-02"
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
	utils.DecodeJSON(resp.Body, &pictures)
	resp.Body.Close()
	return &domain.MainPage{Museums: museums, Exhibitions: exhibitions, Pictures: pictures}
}

func (u *GatewayUsecase) GetPicture(id int) (*domain.Picture, *domain.ErrorResponse) {
	picture := &domain.Picture{}
	resp, err := http.Get(utils.PictureService + strings.Replace(utils.PictureID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		return nil, &domain.ErrorResponse{err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{"Not found"}
	}
	defer resp.Body.Close()
	utils.DecodeJSON(resp.Body, picture)
	return picture, nil
}

func (u *GatewayUsecase) GetExhibition(id int) (*domain.Exhibition, *domain.ErrorResponse) {
	exhibition := &domain.Exhibition{}
	resp, err := http.Get(utils.ExhibitionService + strings.Replace(utils.ExhibitionID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		return nil, &domain.ErrorResponse{err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{"Not found"}
	}
	utils.DecodeJSON(resp.Body, exhibition)
	resp.Body.Close()

	exhibition.Pictures = make([]*domain.Picture, 0)
	resp, err = http.Get(utils.PictureService + utils.PictureByExhibition + fmt.Sprint(exhibition.ID))
	if err != nil {
		return nil, &domain.ErrorResponse{err.Error()}
	} else if resp.StatusCode == http.StatusOK {
		utils.DecodeJSON(resp.Body, &exhibition.Pictures)
	}
	return exhibition, nil
}

func (u *GatewayUsecase) GetMuseum(id int) (*domain.Museum, *domain.ErrorResponse) {
	museum := &domain.Museum{}
	resp, err := http.Get(utils.MuseumService + strings.Replace(utils.MuseumID, ":id", fmt.Sprint(id), 1))
	if err != nil {
		return nil, &domain.ErrorResponse{err.Error()}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &domain.ErrorResponse{"Not found"}
	}
	utils.DecodeJSON(resp.Body, museum)
	resp.Body.Close()

	museum.Exhibitions = make([]*domain.Exhibition, 0)
	resp, err = http.Get(utils.ExhibitionService + utils.ExhibitionByMuseum + fmt.Sprint(museum.ID))
	if err != nil {
		return nil, &domain.ErrorResponse{err.Error()}
	} else if resp.StatusCode == http.StatusOK {
		utils.DecodeJSON(resp.Body, &museum.Exhibitions)
	}
	return museum, nil
}
