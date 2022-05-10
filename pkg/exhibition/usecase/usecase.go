package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/exhibition/repository"
	"log"
)

type ExhibitionUsecase struct {
	repo *repository.ExhibitionRepository
}

func ExhibitionUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.ExhibitionRepository)
	if ok {
		return &ExhibitionUsecase{repo: instance}
	}
	log.Println("Unknown object instead of exhibition usecase")
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

func (u *ExhibitionUsecase) GetExhibitionsByUser(user int) []*domain.Exhibition {
	return u.repo.ExhibitionsByUser(user)
}

func (u *ExhibitionUsecase) GetExhibitionID(id, user int) (*domain.Exhibition, error) {
	if user == 0 {
		exh, err := u.repo.ExhibitionID(id)
		if err == nil {
			u.repo.UpdateExhibitionPopular(id)
		}
		return exh, err
	}
	return u.repo.ExhibitionIDUser(id, user)
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

func (u *ExhibitionUsecase) ShowID(id, user int) error {
	return u.repo.ShowID(id, user)
}

func (u *ExhibitionUsecase) Create(exhibition *domain.Exhibition, museum *domain.Museum, user int) *domain.Exhibition {
	return u.repo.Create(exhibition, museum, user)
}

func (u *ExhibitionUsecase) Update(exhibition *domain.Exhibition, user int) *domain.Exhibition {
	return u.repo.Update(exhibition, user)
}

func (u *ExhibitionUsecase) UpdateImage(exhibition *domain.Exhibition, user int) *domain.Exhibition {
	return u.repo.UpdateImage(exhibition, user)
}

func (u *ExhibitionUsecase) Delete(id, user int) error {
	return u.repo.Delete(id, user)
}
