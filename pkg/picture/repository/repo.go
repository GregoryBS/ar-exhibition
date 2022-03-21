package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"fmt"
)

const (
	querySelectAll   = `select id, name, image, height, width from picture;`
	querySelectByExh = `select id, name, technique, image, author, year, height, width 
	from picture where exh_id = $1;`
)

type PictureRepository struct {
	db *database.DBManager
}

func PictureRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &PictureRepository{db: instance}
	}
	return nil
}

func (repo *PictureRepository) ExhibitionPictures(exhibition int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByExh, exhibition)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Picture{}
		err = rows.Scan(&row.ID, &row.Name, &row.Technique, &row.Image, &row.Author, &row.Year, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		result = append(result, row)
	}
	return result
}

func (repo *PictureRepository) AllPictures() []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectAll)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Picture{}
		err = rows.Scan(&row.ID, &row.Name, &row.Technique, &row.Image, &row.Author, &row.Year, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		result = append(result, row)
	}
	return result
}
