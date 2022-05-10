package repository

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"log"
	"time"
)

const (
	querySelectSearch = `select id, name, image, image_height, image_width, info
	from exhibition where lower(name) like lower($1) and exh_show and mus_show;`
	querySelectSearchID = `select id, name, image, image_height, image_width, info
	from exhibition where lower(name) like lower($1) and museum_id = $2 and exh_show and mus_show;`
)

func (repo *ExhibitionRepository) Search(name, filter string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearch, "%"+name+"%")
	if err != nil {
		log.Println("Cannot search exhibitions by name:", name, err)
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
		row.Image = utils.SplitPic(row.Image)[0]
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
		log.Println("Cannot search exhibitions by name:", name, err)
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
		row.Image = utils.SplitPic(row.Image)[0]
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
