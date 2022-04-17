package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
)

const (
	queryInsert      = `insert into user(login, password) values($1,$2) returning id, login;`
	querySelectLogin = `select id, login, password from user where login = $1;`
)

type UserRepository struct {
	db *database.DBManager
}

func UserRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &UserRepository{db: instance}
	}
	return nil
}

func (repo *UserRepository) Create(user *domain.User) (*domain.User, error) {
	result := &domain.User{}
	row := repo.db.Pool.QueryRow(context.Background(), queryInsert, user.Login, user.Password)
	err := row.Scan(&result.ID, &result.Login)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *UserRepository) GetByLogin(login string) (*domain.User, error) {
	result := &domain.User{}
	row := repo.db.Pool.QueryRow(context.Background(), querySelectLogin, login)
	err := row.Scan(&result.ID, &result.Login, &result.Password)
	if err != nil {
		return nil, err
	}
	return result, nil
}
