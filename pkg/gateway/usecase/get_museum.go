package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
