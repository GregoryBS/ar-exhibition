package usecase

import (
	"ar_exhibition/pkg/user/repository"
)

type UserUsecase struct {
	repo *repository.UserRepository
}

func UserUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.UserRepository)
	if ok {
		return &UserUsecase{repo: instance}
	}
	return nil
}
