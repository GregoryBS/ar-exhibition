package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
)

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
