package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"errors"
	"fmt"
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
	queryUpdatePopular = `update museum set popular = popular + 1 where id = $1;`
	queryInsert        = `insert into museum (name, user_id) values($1, $2) returning id;`
	queryUpdate        = `update museum set name = $1, description = $2, info = $3 where id = $4 and user_id = $5;`
	queryUpdateImage   = `update museum set image = $1, image_height = $2, image_width = $3 where id = $4 and user_id = $5;`
	queryShow          = `update museum set mus_show = not mus_show where id = $1 and user_id = $2;`
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
		return nil, err
	}
	museum.Image = strings.Join(utils.SplitPic(museum.Image), ",")
	museum.Info = utils.MapJSON(params)
	return museum, nil
}

func (repo *MuseumRepository) UpdateMuseumPopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		fmt.Println("Cannot update popular with museum id: ", id)
	}
}

func (repo *MuseumRepository) Museums(page, size int) *domain.Page {
	//offset, limit := (page-1)*size, size
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByPage) //, offset, limit)
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
		row.Image = utils.SplitPic(row.Image)[0]
		result = append(result, row)
	}
	return &domain.Page{Number: 1, Size: len(result), Total: len(result), Items: result}
	//return &domain.Page{Number: page, Size: size, Total: len(result), Items: result}
}

func (repo *MuseumRepository) UserMuseums(user int) []*domain.Museum {
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByUser, user)
	if err != nil {
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

func (repo *MuseumRepository) Create(museum *domain.Museum, user int) *domain.Museum {
	row := repo.db.Pool.QueryRow(context.Background(), queryInsert, museum.Name, user)
	err := row.Scan(&museum.ID)
	if err != nil {
		return nil
	}
	return museum
}

func (repo *MuseumRepository) Update(museum *domain.Museum, user int) *domain.Museum {
	params := make(map[string]string, 0)
	for _, v := range museum.Info {
		params[v.Type] = v.Value
	}
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdate, museum.Name, museum.Description, params, museum.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		return nil
	}
	return museum
}

func (repo *MuseumRepository) UpdateImage(museum *domain.Museum, user int) *domain.Museum {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateImage, museum.Image, museum.Sizes.Height, museum.Sizes.Width, museum.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		return nil
	}
	return museum
}

func (repo *MuseumRepository) Show(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryShow, id, user)
	if err != nil {
		return err
	} else if result.RowsAffected() == 0 {
		return errors.New("Museum not found")
	}
	return nil
}
