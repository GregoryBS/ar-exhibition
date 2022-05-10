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
