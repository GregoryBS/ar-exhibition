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

func (u *PictureUsecase) GetPicturesByExh(exhibition int) []*domain.Picture {
	if exhibition > 0 {
		return u.repo.ExhibitionPictures(exhibition)
	} else {
		return u.repo.TopPictures(15)
	}
}

func (u *PictureUsecase) GetPictureID(id int) (*domain.Picture, error) {
	pic, err := u.repo.PictureID(id)
	if err == nil {
		u.repo.UpdatePicturePopular(id)
	}
	return pic, err
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

func (u *PictureUsecase) Update(pic *domain.Picture, user int) *domain.Picture {
	return u.repo.Update(pic, user)
}
