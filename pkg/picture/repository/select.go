package repository

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	querySelectTop = `select id, name, image, height, width 
	from picture where pic_show and '1' = any (exh_show) and mus_show order by popular desc limit $1;`
	querySelectByExh = `select id, name, image, height, width, video, video_size
	from picture where $1 = any (exh_id)%s;`
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
)

func (repo *PictureRepository) UserPictures(user int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectByUser, user)
	if err != nil {
		log.Println("Cannot get user pictures:", user, err)
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

func (repo *PictureRepository) ExhibitionPictures(exhibition, user int) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	var query string
	if user > 0 {
		query = fmt.Sprintf(querySelectByExh, "")
	} else {
		query = fmt.Sprintf(querySelectByExh, " and pic_show")
	}
	rows, err := repo.db.Pool.Query(context.Background(), query, exhibition)
	if err != nil {
		log.Println("Cannot get exhibition pictures:", exhibition, err)
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
		log.Println("Cannot get top pictures:", err)
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
		log.Println("Cannot get picture:", id, err)
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
		log.Println("Cannot get user picture:", user, id, err)
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

func (repo *PictureRepository) Search(name string) []*domain.Picture {
	result := make([]*domain.Picture, 0)
	rows, err := repo.db.Pool.Query(context.Background(), querySelectSearch, "%"+name+"%")
	if err != nil {
		log.Println("Cannot search pictures with name:", name, err)
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
		log.Println("Cannot search pictures with name:", name, err)
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
