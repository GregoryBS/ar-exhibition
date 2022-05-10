package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"errors"
	"log"
)

const (
	queryUpdatePopular = `update exhibition set popular = popular + 1 where id = $1;`
	queryShow          = `update exhibition set mus_show = not mus_show where user_id = $1;`
	queryShowID        = `update exhibition set exh_show = not exh_show where id = $1 and user_id = $2;`
	queryInsert        = `insert into exhibition (name, description, info, museum_id, exh_show, mus_show, user_id) 
	values($1, $2, $3, $4, $5, $6, $7) returning id;`
	queryUpdate      = `update exhibition set name = $1, description = $2, info = $3 where id = $4 and user_id = $5;`
	queryUpdateImage = `update exhibition set image = $1, image_height = $2, image_width = $3 where id = $4 and user_id = $5;`
	queryDeleteID    = `delete from exhibition where id = $1 and user_id = $2;`
)

const (
	timeLayout = "2006-01-02"
)

type ExhibitionRepository struct {
	db *database.DBManager
}

func ExhibitionRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &ExhibitionRepository{db: instance}
	}
	log.Println("Unknown object instead of db-manager")
	return nil
}

func (repo *ExhibitionRepository) UpdateExhibitionPopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		log.Println("Cannot update popular with exhibition id: ", id, err)
	}
}

func (repo *ExhibitionRepository) Show(user int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryShow, user)
	if err != nil {
		log.Println("User exhibitions publish error:", user, err)
		return err
	}
	return nil
}

func (repo *ExhibitionRepository) ShowID(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryShowID, id, user)
	if err != nil {
		log.Println("Exhibition publish error:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Exhibition for publishing not found")
		return errors.New("Exhibition not found")
	}
	return nil
}

func (repo *ExhibitionRepository) Create(exhibition *domain.Exhibition, museum *domain.Museum, user int) *domain.Exhibition {
	params := make(map[string]string, 0)
	for _, v := range exhibition.Info {
		params[v.Type] = v.Value
	}
	exh_show, mus_show := false, false
	if exhibition.Show > 0 {
		exh_show = true
	}
	if museum.Show > 0 {
		mus_show = true
	}
	row := repo.db.Pool.QueryRow(context.Background(), queryInsert, exhibition.Name, exhibition.Description, params, museum.ID, exh_show, mus_show, user)
	err := row.Scan(&exhibition.ID)
	if err != nil {
		log.Println("Exhibition creating error:", err)
		return nil
	}
	return exhibition
}

func (repo *ExhibitionRepository) Update(exhibition *domain.Exhibition, user int) *domain.Exhibition {
	params := make(map[string]string, 0)
	for _, v := range exhibition.Info {
		params[v.Type] = v.Value
	}
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdate, exhibition.Name, exhibition.Description, params, exhibition.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Exhibition updating error:", exhibition.ID, err, result.RowsAffected())
		return nil
	}
	return exhibition
}

func (repo *ExhibitionRepository) UpdateImage(exhibition *domain.Exhibition, user int) *domain.Exhibition {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateImage,
		exhibition.Image, exhibition.Sizes.Height, exhibition.Sizes.Width, exhibition.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Exhibition image updating error:", exhibition.ID, err, result.RowsAffected())
		return nil
	}
	return exhibition
}

func (repo *ExhibitionRepository) Delete(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryDeleteID, id, user)
	if err != nil {
		log.Println("Exhibition deleting error:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Exhibition for deleting not found")
		return errors.New("Exhibition not found")
	}
	return nil
}
