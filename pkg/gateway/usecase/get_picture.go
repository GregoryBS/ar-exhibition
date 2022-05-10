package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
