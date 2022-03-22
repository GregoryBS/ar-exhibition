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
	from exhibition order by popular desc limit $1;`
	querySelectOne = `select id, name, image, description, info, image_height, image_width
	from exhibition where id = $1;`
	querySelectAll = `select id, name, image, image_height, image_width
	from exhibition;`
	querySelectByMuseum = `select id, name, image, image_height, image_width
	from exhibition where museum_id = $1;`
)

type ExhibitionRepository struct {
	db *database.DBManager
}

func ExhibitionRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &ExhibitionRepository{db: instance}
	}
	return nil
}

func (repo *ExhibitionRepository) ExhibitionTop(limit int) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectTop, limit)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.Service + row.Image
		result = append(result, row)
	}
	return result
}

func (repo *ExhibitionRepository) ExhibitionID(id int) (*domain.Exhibition, error) {
	exh := &domain.Exhibition{Sizes: &domain.ImageSize{}}
	params := make(map[string]string, 0)
	row := repo.db.Pool.QueryRow(context.Background(), querySelectOne, id)
	err := row.Scan(&exh.ID, &exh.Name, &exh.Image, &exh.Description, &params, &exh.Sizes.Height, &exh.Sizes.Width)
	if err != nil {
		return nil, err
	}
	exh.Image = utils.Service + exh.Image
	exh.Info = utils.MapJSON(params)
	return exh, nil
}

func (repo *ExhibitionRepository) ExhibitionByMuseum(museum int) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByMuseum, museum)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.Service + row.Image
		result = append(result, row)
	}
	return result
}

func (repo *ExhibitionRepository) AllExhibitions() []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectAll)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.Service + row.Image
		result = append(result, row)
	}
	return result
}
