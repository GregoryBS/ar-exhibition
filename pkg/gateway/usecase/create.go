package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
)

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
