package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/exhibition/repository"
)

type ExhibitionUsecase struct {
	repo *repository.ExhibitionRepository
}

func ExhibitionUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.ExhibitionRepository)
	if ok {
		return &ExhibitionUsecase{repo: instance}
	}
	return nil
}

func (u *ExhibitionUsecase) GetExhibitionTop() []*domain.Exhibition {
	return u.repo.ExhibitionTop(5)
}
