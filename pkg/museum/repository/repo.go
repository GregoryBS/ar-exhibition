package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"fmt"
)

const querySelectTop = `select id, name, country, city, address, year, description, director, image
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
		err = rows.Scan(&row.ID, &row.Name, &row.Country, &row.City, &row.Address, &row.Year, &row.Description, &row.Director, &row.Image)
		if err != nil {
			return result
		}
		result = append(result, row)
	}
	return result
}
