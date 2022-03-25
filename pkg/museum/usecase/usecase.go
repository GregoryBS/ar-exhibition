package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/museum/repository"
)

type MuseumUsecase struct {
	repo *repository.MuseumRepository
}

func MuseumUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.MuseumRepository)
	if ok {
		return &MuseumUsecase{repo: instance}
	}
	return nil
}

func (u *MuseumUsecase) GetMuseumTop() []*domain.Museum {
	return u.repo.MuseumTop(5)
}

func (u *MuseumUsecase) GetMuseumID(id int) (*domain.Museum, error) {
	return u.repo.MuseumID(id)
}
