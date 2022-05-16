package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/user/repository"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo *repository.UserRepository
}

func UserUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.UserRepository)
	if ok {
		return &UserUsecase{repo: instance}
	}
	log.Println("Unknown object instead of user usecase")
	return nil
}

func (u *UserUsecase) Signup(user *domain.User) (*domain.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.PassIn), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = password
	return u.repo.Create(user)

}

func (u *UserUsecase) Login(user *domain.User) (*domain.User, error) {
	loginUser, err := u.repo.GetByLogin(user.Login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(loginUser.Password, []byte(user.PassIn))
	if err != nil {
		fmt.Println("Passwords do not match")
		return nil, err
	}
	return loginUser, nil
}

func (u UserUsecase) CheckAdmin(id int) error {
	return u.repo.IsAdmin(id)
}
