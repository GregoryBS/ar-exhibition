package usecase

import (
	"ar_exhibition/pkg/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

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
	u.q.Push(req, time.Now().Add(time.Second))
	// resp, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Println("Museum exhibitions publishing error:", museum, err)
	// 	return err
	// }
	// resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return errors.New("Unable to publish museum exhibitions")
	// }

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PictureShow, nil)
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	u.q.Push(req, time.Now().Add(time.Second))
	// resp, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Println("Museum pictures publishing error:", museum, err)
	// 	return err
	// }
	// resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return errors.New("Unable to publish museum pictures")
	// }
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
	u.q.Push(req, time.Now().Add(time.Second))
	// resp, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Println("Exhibition pictures publishing error:", exhibition, err)
	// 	return err
	// }
	// resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return errors.New("Unable to publish exhibition pictures")
	// }
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
