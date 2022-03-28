package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/picture/repository"
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

func (u *PictureUsecase) GetPictures(exhibition int) []*domain.Picture {
	if exhibition > 0 {
		return u.repo.ExhibitionPictures(exhibition)
	} else {
		return u.repo.AllPictures()
	}
}

func (u *PictureUsecase) GetPictureID(id int) (*domain.Picture, error) {
	return u.repo.PictureID(id)
}

func (u *PictureUsecase) Search(name string) []*domain.Picture {
	return u.repo.Search(name)
}
