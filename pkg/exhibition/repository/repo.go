package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"fmt"
)

const querySelectTop = `select id, name, description, date_from, date_to, image
	from exhibition order by popular desc limit $1;`

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
		row := &domain.Exhibition{}
		err = rows.Scan(&row.ID, &row.Name, &row.Description, &row.From, &row.To, &row.Image)
		if err != nil {
			return result
		}
		result = append(result, row)
	}
	return result
}
