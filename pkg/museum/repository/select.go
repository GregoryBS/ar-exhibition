package repository

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"log"
	"strings"
)

const (
	querySelectTop = `select id, name, image, image_height, image_width
	from museum where mus_show order by popular desc limit $1;`
	querySelectOne = `select id, name, image, description, info, image_height, image_width
	from museum where id = $1 and mus_show;`
	querySelectByPage = `select id, name, image, image_height, image_width
	from museum where mus_show;`  // offset $1 limit $2;`
	querySelectByUser = `select id, name, image, description, info, image_height, image_width, mus_show
	from museum where user_id = $1;`
	querySelectSearch = `select id, name, image, image_height, image_width
	from museum where lower(name) like lower($1) and mus_show;`
)

func (repo *MuseumRepository) MuseumTop(limit int) []*domain.Museum {
	result := make([]*domain.Museum, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectTop, limit)
	if err != nil {
		log.Println("Museum top error:", err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Museum{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
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
		log.Println("Museum id error:", id, err)
		return nil, err
	}
	museum.Image = strings.Join(utils.SplitPic(museum.Image), ",")
	museum.Info = utils.MapJSON(params)
	return museum, nil
}

func (repo *MuseumRepository) Museums(page, size int) *domain.Page {
	//offset, limit := (page-1)*size, size
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByPage) //, offset, limit)
	if err != nil {
		log.Println("Museums error:", err)
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
		row.Image = utils.SplitPic(row.Image)[0]
		result = append(result, row)
	}
	return &domain.Page{Number: 1, Size: len(result), Total: len(result), Items: result}
	//return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
}

func (repo *MuseumRepository) UserMuseums(user int) []*domain.Museum {
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByUser, user)
	if err != nil {
		log.Println("User museums error:", user, err)
		return nil
	}
	defer rows.Close()

	result := make([]*domain.Museum, 0)
	for rows.Next() {
		row := &domain.Museum{Sizes: &domain.ImageSize{}}
		params := make(map[string]string, 0)
		flag := false
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Description, &params, &row.Sizes.Height, &row.Sizes.Width, &flag)
		if err != nil {
			return nil
		}
		if flag {
			row.Show = 1
		} else {
			row.Show = -1
		}
		row.Image = strings.Join(utils.SplitPic(row.Image), ",")
		row.Info = utils.MapJSON(params)
		result = append(result, row)
	}
	return result
}

func (repo *MuseumRepository) Search(name string) []*domain.Museum {
	result := make([]*domain.Museum, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearch, "%"+name+"%")
	if err != nil {
		log.Println("Museum search error:", name, err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Museum{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
		result = append(result, row)
	}
	return result
}
