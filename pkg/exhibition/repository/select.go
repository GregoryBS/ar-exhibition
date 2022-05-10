package repository

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"log"
	"strings"
	"time"
)

const (
	querySelectTop = `select id, name, image, image_height, image_width, info
	from exhibition where exh_show and mus_show order by popular desc;`
	querySelectOne = `select id, name, image, description, info, image_height, image_width
	from exhibition where id = $1 and exh_show and mus_show;`
	querySelectOneUser = `select id, name, image, description, info, image_height, image_width, exh_show
	from exhibition where id = $1 and user_id = $2;`
	querySelectAll = `select id, name, image, image_height, image_width, info
	from exhibition where exh_show and mus_show;`  // offset $1 limit $2;`
	querySelectByMuseum = `select id, name, image, image_height, image_width, info
	from exhibition where museum_id = $1 and exh_show and mus_show;`
	querySelectByUser = `select id, name, image, image_height, image_width
	from exhibition where user_id = $1;`
)

func (repo *ExhibitionRepository) ExhibitionTop(limit int) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectTop)
	if err != nil {
		log.Println("Cannot get top exhibitions with error:", err)
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
		row.Image = utils.SplitPic(row.Image)[0]
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

func (repo *ExhibitionRepository) ExhibitionsByUser(user int) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByUser, user)
	if err != nil {
		log.Println("Cannot get user exhibitions:", user, err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Exhibition{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
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
		log.Println("Cannot get exhibition:", id, err)
		return nil, err
	}
	exh.Image = strings.Join(utils.SplitPic(exh.Image), ",")
	exh.Info = utils.MapJSON(params)
	return exh, nil
}

func (repo *ExhibitionRepository) ExhibitionIDUser(id, user int) (*domain.Exhibition, error) {
	exh := &domain.Exhibition{Sizes: &domain.ImageSize{}}
	params := make(map[string]string, 0)
	flag := false
	row := repo.db.Pool.QueryRow(context.Background(), querySelectOneUser, id, user)
	err := row.Scan(&exh.ID, &exh.Name, &exh.Image, &exh.Description, &params, &exh.Sizes.Height, &exh.Sizes.Width, &flag)
	if err != nil {
		log.Println("Cannot get user exhibition:", user, id, err)
		return nil, err
	}
	if flag {
		exh.Show = 1
	} else {
		exh.Show = -1
	}
	exh.Image = strings.Join(utils.SplitPic(exh.Image), ",")
	exh.Info = utils.MapJSON(params)
	return exh, nil
}

func (repo *ExhibitionRepository) ExhibitionByMuseum(museum int, filter string) []*domain.Exhibition {
	result := make([]*domain.Exhibition, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByMuseum, museum)
	if err != nil {
		log.Println("Cannot get museum exhibitions:", museum, err)
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

func (repo *ExhibitionRepository) AllExhibitions(page, size int, filter string) *domain.Page {
	//offset, limit := (page-1)*size, size
	rows, err := repo.db.Pool.Query(context.Background(), querySelectAll) //, offset, limit)
	if err != nil {
		log.Println("Cannot get all exhibitions:", err)
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
	return &domain.Page{Number: 1, Size: len(result), Total: len(result), Items: result}
	//return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
}
