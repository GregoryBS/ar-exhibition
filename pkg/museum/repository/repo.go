package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"errors"
	"log"
)

const (
	queryUpdatePopular = `update museum set popular = popular + 1 where id = $1;`
	queryInsert        = `insert into museum (name, user_id) values($1, $2) returning id;`
	queryUpdate        = `update museum set name = $1, description = $2, info = $3 where id = $4 and user_id = $5;`
	queryUpdateImage   = `update museum set image = $1, image_height = $2, image_width = $3 where id = $4 and user_id = $5;`
	queryShow          = `update museum set mus_show = not mus_show where id = $1 and user_id = $2;`
)

type MuseumRepository struct {
	db *database.DBManager
}

func MuseumRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &MuseumRepository{db: instance}
	}
	log.Println("Unknown object instead of db-manager")
	return nil
}

func (repo *MuseumRepository) UpdateMuseumPopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		log.Println("Cannot update popular with museum id: ", id, err)
	}
}

func (repo *MuseumRepository) Create(museum *domain.Museum, user int) *domain.Museum {
	row := repo.db.Pool.QueryRow(context.Background(), queryInsert, museum.Name, user)
	err := row.Scan(&museum.ID)
	if err != nil {
		log.Println("Museum creating error:", err)
		return nil
	}
	return museum
}

func (repo *MuseumRepository) Update(museum *domain.Museum, user int) *domain.Museum {
	params := make(map[string]string, 0)
	for _, v := range museum.Info {
		params[v.Type] = v.Value
	}
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdate, museum.Name, museum.Description, params, museum.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Museum update error:", museum.ID, err, result.RowsAffected())
		return nil
	}
	return museum
}

func (repo *MuseumRepository) UpdateImage(museum *domain.Museum, user int) *domain.Museum {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateImage, museum.Image, museum.Sizes.Height, museum.Sizes.Width, museum.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Museum image update error:", museum.ID, err, result.RowsAffected())
		return nil
	}
	return museum
}

func (repo *MuseumRepository) Show(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryShow, id, user)
	if err != nil {
		log.Println("Museum publish error:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Museum for publishing not found")
		return errors.New("Museum not found")
	}
	return nil
}
