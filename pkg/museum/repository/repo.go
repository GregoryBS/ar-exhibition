package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"fmt"
)

const (
	querySelectTop = `select id, name, image, image_height, image_width
	from museum order by popular desc limit $1;`
	querySelectOne = `select id, name, image, description, info, image_height, image_width
	from museum where id = $1;`
	querySelectByPage = `select id, name, image, image_height, image_width
	from museum offset $1 limit $2;`
	querySelectSearch = `select id, name, image, image_height, image_width
	from museum where lower(name) like lower($1);`
)

type MuseumRepository struct {
	db *database.DBManager
}

func MuseumRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &MuseumRepository{db: instance}
	}
	return nil
}

func (repo *MuseumRepository) MuseumTop(limit int) []*domain.Museum {
	result := make([]*domain.Museum, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectTop, limit)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Museum{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.ImageService + row.Image
		result = append(result, row)
	}
	return result
}

func (repo *MuseumRepository) MuseumID(id int) (*domain.Museum, error) {
	museum := &domain.Museum{Sizes: &domain.ImageSize{}}
	params := make(map[string]string, 0)
	row := repo.db.Pool.QueryRow(context.Background(), querySelectOne, id)
	err := row.Scan(&museum.ID, &museum.Name, &museum.Image, &museum.Description, &params, &museum.Sizes.Height, &museum.Sizes.Width)
	if err != nil {
		return nil, err
	}
	museum.Image = utils.ImageService + museum.Image
	museum.Info = utils.MapJSON(params)
	return museum, nil
}

func (repo *MuseumRepository) Museums(page, size int) *domain.Page {
	offset, limit := (page-1)*size, size
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByPage, offset, limit)
	if err != nil {
		return &domain.Page{Number: page, Size: size}
	}
	defer rows.Close()

	result := make([]interface{}, 0)
	for rows.Next() {
		row := &domain.Museum{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
		}
		row.Image = utils.ImageService + row.Image
		result = append(result, row)
	}
	return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
}

func (repo *MuseumRepository) Search(name string) []*domain.Museum {
	result := make([]*domain.Museum, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearch, "%"+name+"%")
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Museum{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.ImageService + row.Image
		result = append(result, row)
	}
	return result
}
