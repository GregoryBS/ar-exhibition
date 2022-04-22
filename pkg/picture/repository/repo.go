package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"fmt"
	"strings"
)

const (
	querySelectTop = `select id, name, image, height, width 
	from picture where pic_show and exh_show and mus_show order by popular desc limit $1;`
	querySelectByExh = `select id, name, image, height, width, video, video_size
	from picture where exh_id = $1 and pic_show and exh_show and mus_show;`
	querySelectOne = `select id, name, image, description, info, height, width
	from picture where id = $1 and pic_show and exh_show and mus_show;`
	querySelectSearch = `select  id, name, image, height, width 
	from picture where lower(name) like lower($1) and pic_show and exh_show and mus_show;`
	querySelectSearchID = `select  id, name, image, height, width 
	from picture where lower(name) like lower($1) and exh_id = $2 and pic_show and exh_show and mus_show;`
	queryUpdatePopular = `update picture set popular = popular + 1 where id = $1;`
	queryInsert        = `insert into picture (name, description, info, height, width, user_id) values($1, $2, $3, $4, $5, $6) returning id;`
	queryUpdate        = `update picture set name = $1, description = $2, info = $3, height = $6, width = $7 where id = $4 and user_id = $5;`
	queryUpdateImage   = `update picture set image = $1 where id = $2 and user_id = $3;`
	queryShow          = `update picture set mus_show = not mus_show where user_id = $1;`
	queryShowExh       = `update picture set exh_show = not exh_show where exh_id = $1 and user_id = $2;`
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
		row := &domain.Picture{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width, &row.Video, &row.VideoSize)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
		if row.Video != "" {
			row.Video = utils.VideoService + row.Video
		}
		result = append(result, row)
	}
	return result
}

func (repo *PictureRepository) TopPictures(limit int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectTop, limit)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Picture{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
		result = append(result, row)
	}
	return result
}

func (repo *PictureRepository) PictureID(id int) (*domain.Picture, error) {
	pic := &domain.Picture{Sizes: &domain.ImageSize{}}
	params := make(map[string]string, 0)
	row := repo.db.Pool.QueryRow(context.Background(), querySelectOne, id)
	err := row.Scan(&pic.ID, &pic.Name, &pic.Image, &pic.Description, &params, &pic.Sizes.Height, &pic.Sizes.Width)
	if err != nil {
		return nil, err
	}
	pic.Image = strings.Join(utils.SplitPic(pic.Image), ",")
	pic.Info = utils.MapJSON(params)
	return pic, nil
}

func (repo *PictureRepository) UpdatePicturePopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		fmt.Println("Cannot update popular with picture id: ", id)
	}
}

func (repo *PictureRepository) Search(name string) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearch, "%"+name+"%")
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Picture{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
		result = append(result, row)
	}
	return result
}

func (repo *PictureRepository) SearchID(name string, exhibitionID int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearchID, "%"+name+"%", exhibitionID)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		row := &domain.Picture{Sizes: &domain.ImageSize{}}
		err = rows.Scan(&row.ID, &row.Name, &row.Image, &row.Sizes.Height, &row.Sizes.Width)
		if err != nil {
			return result
		}
		row.Image = utils.SplitPic(row.Image)[0]
		result = append(result, row)
	}
	return result
}

func (repo *PictureRepository) Create(picture *domain.Picture, user int) *domain.Picture {
	params := make(map[string]string, 0)
	for _, v := range picture.Info {
		params[v.Type] = v.Value
	}
	row := repo.db.Pool.QueryRow(context.Background(), queryInsert,
		picture.Name, picture.Description, params, picture.Sizes.Height, picture.Sizes.Width, user)
	err := row.Scan(&picture.ID)
	if err != nil {
		return nil
	}
	return picture
}

func (repo *PictureRepository) Update(picture *domain.Picture, user int) *domain.Picture {
	params := make(map[string]string, 0)
	for _, v := range picture.Info {
		params[v.Type] = v.Value
	}
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdate,
		picture.Name, picture.Description, params, picture.ID, user, picture.Sizes.Height, picture.Sizes.Width)
	if err != nil || result.RowsAffected() == 0 {
		return nil
	}
	return picture
}

func (repo *PictureRepository) UpdateImage(picture *domain.Picture, user int) *domain.Picture {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateImage, picture.Image, picture.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		return nil
	}
	return picture
}

func (repo *PictureRepository) Show(user int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryShow, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PictureRepository) ShowExh(exhibiton, user int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryShowExh, exhibiton, user)
	if err != nil {
		return err
	}
	return nil
}
