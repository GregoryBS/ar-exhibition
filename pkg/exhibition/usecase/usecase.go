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

func (u *ExhibitionUsecase) GetExhibitionsByMuseum(museum int, filter string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	if museum > 0 {
		return u.repo.ExhibitionByMuseum(museum, filter)
	}
	return result
}

func (u *ExhibitionUsecase) GetExhibitionID(id int) (*domain.Exhibition, error) {
	exh, err := u.repo.ExhibitionID(id)
	if err == nil {
		u.repo.UpdateExhibitionPopular(id)
	}
	return exh, err
}

func (u *ExhibitionUsecase) GetExhibitions(page, size int, filter string) *domain.Page {
	return u.repo.AllExhibitions(page, size, filter)
}

func (u *ExhibitionUsecase) Search(name, filter string) []*domain.Exhibition {
	return u.repo.Search(name, filter)
}

func (u *ExhibitionUsecase) SearchID(name string, museumID int, filter string) []*domain.Exhibition {
	return u.repo.SearchID(name, museumID, filter)
}

func (u *ExhibitionUsecase) Show(user int) error {
	return u.repo.Show(user)
}
