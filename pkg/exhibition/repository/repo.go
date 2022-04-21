package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"fmt"
	"time"
)

const (
	querySelectTop = `select id, name, image, image_height, image_width, info
	from exhibition where exh_show and mus_show order by popular desc;`
	querySelectOne = `select id, name, image, description, info, image_height, image_width
	from exhibition where id = $1 and exh_show and mus_show;`
	querySelectAll = `select id, name, image, image_height, image_width, info
	from exhibition where and exh_show and mus_show offset $1 limit $2;`
	querySelectByMuseum = `select id, name, image, image_height, image_width, info
	from exhibition where museum_id = $1 and exh_show and mus_show;`
	querySelectSearch = `select id, name, image, image_height, image_width, info
	from exhibition where lower(name) like lower($1) and exh_show and mus_show;`
	querySelectSearchID = `select id, name, image, image_height, image_width, info
	from exhibition where lower(name) like lower($1) and museum_id = $2 and exh_show and mus_show;`
	queryUpdatePopular = `update exhibition set popular = popular + 1 where id = $1;`
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
	return nil
}

func (repo *ExhibitionRepository) ExhibitionTop(limit int) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectTop)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() && limit > 0 {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		params := make(map[string]string, 0)
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width, &params)
		if err != nil {
			return result
		}
		row.Image = utils.ImageService + row.Image
		from, _ := time.Parse(timeLayout, params[utils.ExhibitionStart])
		to, _ := time.Parse(timeLayout, params[utils.ExhibitionEnd])
		t := time.Now()
		if t.Before(to) && t.After(from) {
			result = append(result, row)
			limit -= 1
		}
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
	exh.Image = utils.ImageService + exh.Image
	exh.Info = utils.MapJSON(params)
	return exh, nil
}

func (repo *ExhibitionRepository) UpdateExhibitionPopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		fmt.Println("Cannot update popular with exhibition id: ", id)
	}
}

func (repo *ExhibitionRepository) ExhibitionByMuseum(museum int, filter string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByMuseum, museum)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		params := make(map[string]string, 0)
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width, &params)
		if err != nil {
			return result
		}
		row.Image = utils.ImageService + row.Image
		from, _ := time.Parse(timeLayout, params[utils.ExhibitionStart])
		to, _ := time.Parse(timeLayout, params[utils.ExhibitionEnd])
		switch filter {
		case "all":
			result = append(result, row)
		case "now":
			t := time.Now()
			if t.Before(to) && t.After(from) {
				result = append(result, row)
			}
		case "old":
			if time.Now().After(to) {
				result = append(result, row)
			}
		default:
			result = append(result, row)
		}
	}
	return result
}

func (repo *ExhibitionRepository) AllExhibitions(page, size int, filter string) *domain.Page {
	offset, limit := (page-1)*size, size
	rows, err := repo.db.Pool.Query(context.Background(), querySelectAll, offset, limit)
	if err != nil {
		return &domain.Page{Number: page, Size: size}
	}
	defer rows.Close()

	result := make([]interface{}, 0)
	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		params := make(map[string]string, 0)
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width, &params)
		if err != nil {
			return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
		}
		row.Image = utils.ImageService + row.Image
		from, _ := time.Parse(timeLayout, params[utils.ExhibitionStart])
		to, _ := time.Parse(timeLayout, params[utils.ExhibitionEnd])
		switch filter {
		case "all":
			result = append(result, row)
		case "now":
			t := time.Now()
			if t.Before(to) && t.After(from) {
				result = append(result, row)
			}
		case "old":
			if time.Now().After(to) {
				result = append(result, row)
			}
		default:
			result = append(result, row)
		}
	}
	return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
}

func (repo *ExhibitionRepository) Search(name, filter string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearch, "%"+name+"%")
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		params := make(map[string]string, 0)
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width, &params)
		if err != nil {
			return result
		}
		row.Image = utils.ImageService + row.Image
		from, _ := time.Parse(timeLayout, params[utils.ExhibitionStart])
		to, _ := time.Parse(timeLayout, params[utils.ExhibitionEnd])
		switch filter {
		case "all":
			result = append(result, row)
		case "now":
			t := time.Now()
			if t.Before(to) && t.After(from) {
				result = append(result, row)
			}
		case "old":
			if time.Now().After(to) {
				result = append(result, row)
			}
		default:
			result = append(result, row)
		}
	}
	return result
}

func (repo *ExhibitionRepository) SearchID(name string, museumID int, filter string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearchID, "%"+name+"%", museumID)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		params := make(map[string]string, 0)
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width, &params)
		if err != nil {
			return result
		}
		row.Image = utils.ImageService + row.Image
		from, _ := time.Parse(timeLayout, params[utils.ExhibitionStart])
		to, _ := time.Parse(timeLayout, params[utils.ExhibitionEnd])
		switch filter {
		case "all":
			result = append(result, row)
		case "now":
			t := time.Now()
			if t.Before(to) && t.After(from) {
				result = append(result, row)
			}
		case "old":
			if time.Now().After(to) {
				result = append(result, row)
			}
		default:
			result = append(result, row)
		}
	}
	return result
}
