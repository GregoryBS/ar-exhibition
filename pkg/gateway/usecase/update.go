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
	"time"
)

func (u *GatewayUsecase) UpdateMuseum(museum *domain.Museum, user int) (*domain.Museum, error) {
	req, _ := http.NewRequest(http.MethodPost, utils.MuseumService+strings.Replace(utils.MuseumID, ":id", fmt.Sprint(museum.ID), 1),
		bytes.NewBuffer(utils.EncodeJSON(museum)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum updating error:", err)
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
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Museum image updating error:", err)
		return &domain.ErrorResponse{Message: "Museum service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update museum image"}
	}
	return nil
}

func (u *GatewayUsecase) UpdatePicture(picture *domain.Picture, user int) (*domain.Picture, error) {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureID, ":id", fmt.Sprint(picture.ID), 1),
		bytes.NewBuffer(utils.EncodeJSON(picture)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture updating error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to update picture")
	}
	result := new(domain.Picture)
	utils.DecodeJSON(resp.Body, result)
	return result, nil
}

func (u *GatewayUsecase) UpdatePictureImage(filename string, sizes *domain.ImageSize, pic, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureImage, ":id", fmt.Sprint(pic), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Picture{ID: pic, Image: filename, Sizes: sizes})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture image updating error:", err)
		return &domain.ErrorResponse{Message: "Picture service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update picture image"}
	}
	return nil
}

func (u *GatewayUsecase) UpdatePictureVideo(filename, size string, pic, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.PictureService+strings.Replace(utils.PictureVideo, ":id", fmt.Sprint(pic), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Picture{ID: pic, Video: filename, VideoSize: size})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Picture video updating error:", err)
		return &domain.ErrorResponse{Message: "Picture service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update picture video"}
	}
	return nil
}

func (u *GatewayUsecase) UpdateExhibition(exhibition *domain.Exhibition, user int) (*domain.Exhibition, error) {
	museum := u.GetUserMuseums(user)
	if museum == nil {
		log.Println("Cannot get user museum for exhibition change")
		return nil, errors.New("Unable to update exhibition")
	}

	req, _ := http.NewRequest(http.MethodPost, utils.ExhibitionService+strings.Replace(utils.ExhibitionID, ":id", fmt.Sprint(exhibition.ID), 1),
		bytes.NewBuffer(utils.EncodeJSON(exhibition)))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition updating error:", exhibition.ID, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to update exhibition")
	}
	result := new(domain.Exhibition)
	utils.DecodeJSON(resp.Body, result)
	resp.Body.Close()

	req, _ = http.NewRequest(http.MethodPost, utils.PictureService+utils.PicturesToExh,
		bytes.NewBuffer(utils.EncodeJSON(domain.MuseumExhibition{Mus: museum[0], Exh: result})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	u.q.Push(req, time.Now().Add(time.Millisecond))

	// resp, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Println("Pictures to exhibition adding error:", result.ID, err)
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, errors.New("Unable to update exhibition")
	// }
	return result, nil
}

func (u *GatewayUsecase) UpdateExhibitionImage(filename string, sizes *domain.ImageSize, id, user int) *domain.ErrorResponse {
	req, _ := http.NewRequest(http.MethodPost, utils.ExhibitionService+strings.Replace(utils.ExhibitionImage, ":id", fmt.Sprint(id), 1),
		bytes.NewBuffer(utils.EncodeJSON(domain.Exhibition{ID: id, Image: filename, Sizes: sizes})))
	req.Header.Set(utils.UserHeader, fmt.Sprint(user))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Exhibition image updating error:", id, err)
		return &domain.ErrorResponse{Message: "Exhibition service is unavailable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &domain.ErrorResponse{Message: "Unable to update exhibition image"}
	}
	return nil
}
