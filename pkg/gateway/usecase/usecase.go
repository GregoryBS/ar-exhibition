package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"net/http"
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
