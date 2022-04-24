package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
)

const (
	querySelectTop = `select id, name, image, height, width 
	from picture where pic_show and '1' = any (exh_show) and mus_show order by popular desc limit $1;`
	querySelectByExh = `select id, name, image, height, width, video, video_size
	from picture where $1 = any (exh_id) and pic_show and '1' = any (exh_show) and mus_show;`
	querySelectByUser = `select id, name, image, height, width
	from picture where user_id = $1;`
	querySelectOne = `select id, name, image, description, info, height, width
	from picture where id = $1 and pic_show and '1' = any (exh_show) and mus_show;`
	querySelectOneByUser = `select id, name, image, description, info, height, width, video, video_size, pic_show
	from picture where id = $1 and user_id = $2;`
	querySelectSearch = `select  id, name, image, height, width 
	from picture where lower(name) like lower($1) and pic_show and '1' = any (exh_show) and mus_show;`
	querySelectSearchID = `select  id, name, image, height, width 
	from picture where lower(name) like lower($1) and $2 = any(exh_id) and pic_show and '1' = any (exh_show) and mus_show;`
	queryUpdatePopular = `update picture set popular = popular + 1 where id = $1;`
	queryInsert        = `insert into picture (name, description, info, height, width, user_id) values($1, $2, $3, $4, $5, $6) returning id;`
	queryUpdate        = `update picture set name = $1, description = $2, info = $3, height = $6, width = $7 where id = $4 and user_id = $5;`
	queryUpdateImage   = `update picture set image = $1 where id = $2 and user_id = $3;`
	queryUpdateVideo   = `update picture set video = $1, video_size = $2 where id = $3 and user_id = $4;`
	queryShow          = `update picture set mus_show = not mus_show where user_id = $1;`
	queryShowExh       = `update picture set exh_show[array_position(exh_id, $1)] = not exh_show[array_position(exh_id, $1)] 
	where $1 = any (exh_id) and user_id = $2;`
	queryShowID           = `update picture set pic_show = not pic_show where id = $1 and user_id = $2;`
	queryDeleteID         = `delete from picture where id = $1 and user_id = $2;`
	queryDeleteExhibition = `update picture set exh_id = array_remove(exh_id, $1), 
	exh_show = exh_show[1:array_position(exh_id, $1)-1] || exh_show[array_position(exh_id, $1)+1:] where $1 = any (exh_id);`
	queryAddExhibition = `update picture set exh_id = array_append(exh_id,$1), exh_show = array_append(exh_show,$2)%s where id = $3 and user_id = $4;`
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

func (repo *PictureRepository) UserPictures(user int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByUser, user)
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

func (repo *PictureRepository) PictureIDUser(id, user int) (*domain.Picture, error) {
	pic := &domain.Picture{Sizes: &domain.ImageSize{}}
	params := make(map[string]string, 0)
	flag := false
	row := repo.db.Pool.QueryRow(context.Background(), querySelectOneByUser, id, user)
	err := row.Scan(&pic.ID, &pic.Name, &pic.Image, &pic.Description, &params, &pic.Sizes.Height, &pic.Sizes.Width, &pic.Video, &pic.VideoSize, &flag)
	if err != nil {
		return nil, err
	}
	if flag {
		pic.Show = 1
	} else {
		pic.Show = -1
	}
	pic.Image = strings.Join(utils.SplitPic(pic.Image), ",")
	if pic.Video != "" {
		pic.Video = utils.VideoService + pic.Video
	}
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

func (repo *PictureRepository) UpdateVideo(picture *domain.Picture, user int) *domain.Picture {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateVideo, picture.Video, picture.VideoSize, picture.ID, user)
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

func (repo *PictureRepository) ShowID(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryShowID, id, user)
	if err != nil {
		return err
	} else if result.RowsAffected() == 0 {
		return errors.New("Picture not found")
	}
	return nil
}

func (repo *PictureRepository) Delete(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryDeleteID, id, user)
	if err != nil {
		return err
	} else if result.RowsAffected() == 0 {
		return errors.New("Picture not found")
	}
	return nil
}

func (repo *PictureRepository) DeleteFromExhibition(exhibition int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryDeleteExhibition, exhibition)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PictureRepository) AddToExhibition(pic *domain.Picture, exh *domain.Exhibition, mus *domain.Museum, user int) error {
	exh_flag := false
	if exh.Show > 0 {
		exh_flag = true
	}
	var result pgconn.CommandTag
	var err error
	if mus == nil {
		sql := fmt.Sprintf(queryAddExhibition, "")
		result, err = repo.db.Pool.Exec(context.Background(), sql, exh.ID, exh_flag, pic.ID, user)
	} else {
		mus_flag := false
		if mus.Show > 0 {
			mus_flag = true
		}
		sql := fmt.Sprintf(queryAddExhibition, ", mus_show = $5")
		result, err = repo.db.Pool.Exec(context.Background(), sql, exh.ID, exh_flag, pic.ID, user, mus_flag)
	}
	if err != nil {
		return err
	} else if result.RowsAffected() == 0 {
		return errors.New("Picture not found")
	}
	return nil
}
