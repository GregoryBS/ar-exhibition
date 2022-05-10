package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
