package repository

import "ar_exhibition/pkg/database"

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
