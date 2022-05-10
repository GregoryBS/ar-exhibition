package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
)

const (
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
	log.Println("Unknown object instead of db-manager")
	return nil
}

func (repo *PictureRepository) UpdatePicturePopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		log.Println("Cannot update popular with picture id: ", id, err)
	}
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
		log.Println("Cannot create picture:", err)
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
		log.Println("Cannot update picture:", picture.ID, err, result.RowsAffected())
		return nil
	}
	return picture
}

func (repo *PictureRepository) UpdateImage(picture *domain.Picture, user int) *domain.Picture {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateImage, picture.Image, picture.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Cannot update picture image:", picture.ID, err, result.RowsAffected())
		return nil
	}
	return picture
}

func (repo *PictureRepository) UpdateVideo(picture *domain.Picture, user int) *domain.Picture {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateVideo, picture.Video, picture.VideoSize, picture.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Cannot update picture video:", picture.ID, err, result.RowsAffected())
		return nil
	}
	return picture
}

func (repo *PictureRepository) Show(user int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryShow, user)
	if err != nil {
		log.Println("Cannot publish user pictures:", user, err)
		return err
	}
	return nil
}

func (repo *PictureRepository) ShowExh(exhibition, user int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryShowExh, exhibition, user)
	if err != nil {
		log.Println("Cannot publish exhibition pictures:", exhibition, err)
		return err
	}
	return nil
}

func (repo *PictureRepository) ShowID(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryShowID, id, user)
	if err != nil {
		log.Println("Cannot publish picture:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Picture for publish not found")
		return errors.New("Picture not found")
	}
	return nil
}

func (repo *PictureRepository) Delete(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryDeleteID, id, user)
	if err != nil {
		log.Println("Cannot delete picture:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Picture for deleting not found")
		return errors.New("Picture not found")
	}
	return nil
}

func (repo *PictureRepository) DeleteFromExhibition(exhibition int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryDeleteExhibition, exhibition)
	if err != nil {
		log.Println("Cannot delete pictures from exhibition:", exhibition, err)
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
		log.Println("Cannot add picture to exhibition:", pic.ID, exh.ID, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Cannot find picture to add to exhibition:", pic.ID, exh.ID, err)
		return errors.New("Picture not found")
	}
	return nil
}
