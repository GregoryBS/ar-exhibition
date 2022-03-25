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

func (u *ExhibitionUsecase) GetExhibitionsByMuseum(museum int) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	if museum > 0 {
		return u.repo.ExhibitionByMuseum(museum)
	}
	return result
}

func (u *ExhibitionUsecase) GetExhibitionID(id int) (*domain.Exhibition, error) {
	return u.repo.ExhibitionID(id)
}

func (u *ExhibitionUsecase) GetExhibitions(page, size int) *domain.Page {
	return u.repo.AllExhibitions(page, size)
}
