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

func (u *PictureUsecase) GetPicturesByExh(exhibition int) []*domain.Picture {
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

func (u *PictureUsecase) SearchID(name string, exhibitionID int) []*domain.Picture {
	return u.repo.SearchID(name, exhibitionID)
}

func (u *PictureUsecase) GetPicturesByIDs(id []int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	for i := range id {
		if pic, err := u.GetPictureID(id[i]); err == nil {
			result = append(result, pic)
		}
	}
	return result
}
