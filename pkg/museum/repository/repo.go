package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"fmt"
)

const querySelectTop = `select id, name, image, image_height, image_width
	from museum order by popular desc limit $1;`

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
		row := &domain.Museum{}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.Service + row.Image
		result = append(result, row)
	}
	return result
}
