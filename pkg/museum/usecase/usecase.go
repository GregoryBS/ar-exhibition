package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/museum/repository"
	"log"
)

type MuseumUsecase struct {
	repo *repository.MuseumRepository
}

func MuseumUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.MuseumRepository)
	if ok {
		return &MuseumUsecase{repo: instance}
	}
	log.Println("Unknown object instead of museum usecase")
	return nil
}

func (u *MuseumUsecase) GetMuseumTop() []*domain.Museum {
	return u.repo.MuseumTop(5)
}

func (u *MuseumUsecase) GetMuseumID(id int) (*domain.Museum, error) {
	museum, err := u.repo.MuseumID(id)
	if err == nil {
		u.repo.UpdateMuseumPopular(id)
	}
	return museum, err
}

func (u *MuseumUsecase) GetMuseums(page, size int) *domain.Page {
	return u.repo.Museums(page, size)
}

func (u *MuseumUsecase) GetUserMuseums(user int) []*domain.Museum {
	return u.repo.UserMuseums(user)
}

func (u *MuseumUsecase) Search(name string) []*domain.Museum {
	return u.repo.Search(name)
}

func (u *MuseumUsecase) Create(museum *domain.Museum, user int) *domain.Museum {
	return u.repo.Create(museum, user)
}

func (u *MuseumUsecase) Update(museum *domain.Museum, user int) *domain.Museum {
	return u.repo.Update(museum, user)
}

func (u *MuseumUsecase) UpdateImage(museum *domain.Museum, user int) *domain.Museum {
	return u.repo.UpdateImage(museum, user)
}

func (u *MuseumUsecase) Show(id, user int) error {
	return u.repo.Show(id, user)
}
