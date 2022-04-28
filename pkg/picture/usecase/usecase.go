package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/picture/repository"
	"strings"
)

type PictureUsecase struct {
	repo *repository.PictureRepository
}

func PictureUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.PictureRepository)
	if ok {
		return &PictureUsecase{repo: instance}
	}
	return nil
}

func (u *PictureUsecase) GetPicturesByExh(exhibition, user int) []*domain.Picture {
	if exhibition > 0 {
		return u.repo.ExhibitionPictures(exhibition, user)
	} else {
		return u.repo.TopPictures(15)
	}
}

func (u *PictureUsecase) GetPicturesByUser(user int) []*domain.Picture {
	return u.repo.UserPictures(user)
}

func (u *PictureUsecase) GetPictureID(id, user int) (*domain.Picture, error) {
	if user == 0 {
		pic, err := u.repo.PictureID(id)
		if err == nil {
			u.repo.UpdatePicturePopular(id)
		}
		return pic, err
	}
	return u.repo.PictureIDUser(id, user)
}

func (u *PictureUsecase) Search(name string) []*domain.Picture {
	return u.repo.Search(name)
}

func (u *PictureUsecase) SearchID(name string, exhibitionID int) []*domain.Picture {
	return u.repo.SearchID(name, exhibitionID)
}

func (u *PictureUsecase) GetPicturesByIDs(id []int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	for i := range id {
		if pic, err := u.repo.PictureID(id[i]); err == nil {
			pic.Info, pic.Description = nil, ""
			pic.Image = strings.Split(pic.Image, ",")[0]
			result = append(result, pic)
		}
	}
	return result
}

func (u *PictureUsecase) Create(pic *domain.Picture, user int) *domain.Picture {
	return u.repo.Create(pic, user)
}

func (u *PictureUsecase) UpdateImage(pic *domain.Picture, user int) *domain.Picture {
	return u.repo.UpdateImage(pic, user)
}

func (u *PictureUsecase) UpdateVideo(pic *domain.Picture, user int) *domain.Picture {
	return u.repo.UpdateVideo(pic, user)
}

func (u *PictureUsecase) Update(pic *domain.Picture, user int) *domain.Picture {
	return u.repo.Update(pic, user)
}

func (u *PictureUsecase) Show(user int) error {
	return u.repo.Show(user)
}

func (u *PictureUsecase) ShowExh(exhibition, user int) error {
	return u.repo.ShowExh(exhibition, user)
}

func (u *PictureUsecase) ShowID(id, user int) error {
	return u.repo.ShowID(id, user)
}

func (u *PictureUsecase) Delete(id, user int) error {
	return u.repo.Delete(id, user)
}

func (u *PictureUsecase) UpdateForExhibition(exh *domain.Exhibition, mus *domain.Museum, user int) error {
	err := u.repo.DeleteFromExhibition(exh.ID)
	if err == nil {
		for i := range exh.Pictures {
			err = u.repo.AddToExhibition(exh.Pictures[i], exh, mus, user)
			if err != nil {
				break
			}
		}
	}
	return err
}
