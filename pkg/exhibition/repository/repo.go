package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"context"
	"errors"
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
	querySelectSearch = `select id, name, image, image_height, image_width, info
	from exhibition where lower(name) like lower($1) and exh_show and mus_show;`
	querySelectSearchID = `select id, name, image, image_height, image_width, info
	from exhibition where lower(name) like lower($1) and museum_id = $2 and exh_show and mus_show;`
	queryUpdatePopular = `update exhibition set popular = popular + 1 where id = $1;`
	queryShow          = `update exhibition set mus_show = not mus_show where user_id = $1;`
	queryShowID        = `update exhibition set exh_show = not exh_show where id = $1 and user_id = $2;`
	queryInsert        = `insert into exhibition (name, description, info, museum_id, exh_show, mus_show, user_id) 
	values($1, $2, $3, $4, $5, $6, $7) returning id;`
	queryUpdate      = `update exhibition set name = $1, description = $2, info = $3 where id = $4 and user_id = $5;`
	queryUpdateImage = `update exhibition set image = $1, image_height = $2, image_width = $3 where id = $4 and user_id = $5;`
	queryDeleteID    = `delete from exhibition where id = $1 and user_id = $2;`
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
	log.Println("Unknown object instead of db-manager")
	return nil
}

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

func (repo *ExhibitionRepository) UpdateExhibitionPopular(id int) {
	_, err := repo.db.Pool.Exec(context.Background(), queryUpdatePopular, id)
	if err != nil {
		log.Println("Cannot update popular with exhibition id: ", id, err)
	}
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

func (repo *ExhibitionRepository) Show(user int) error {
	_, err := repo.db.Pool.Exec(context.Background(), queryShow, user)
	if err != nil {
		log.Println("User exhibitions publish error:", user, err)
		return err
	}
	return nil
}

func (repo *ExhibitionRepository) ShowID(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryShowID, id, user)
	if err != nil {
		log.Println("Exhibition publish error:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Exhibition for publishing not found")
		return errors.New("Exhibition not found")
	}
	return nil
}

func (repo *ExhibitionRepository) Create(exhibition *domain.Exhibition, museum *domain.Museum, user int) *domain.Exhibition {
	params := make(map[string]string, 0)
	for _, v := range exhibition.Info {
		params[v.Type] = v.Value
	}
	exh_show, mus_show := false, false
	if exhibition.Show > 0 {
		exh_show = true
	}
	if museum.Show > 0 {
		mus_show = true
	}
	row := repo.db.Pool.QueryRow(context.Background(), queryInsert, exhibition.Name, exhibition.Description, params, museum.ID, exh_show, mus_show, user)
	err := row.Scan(&exhibition.ID)
	if err != nil {
		log.Println("Exhibition creating error:", err)
		return nil
	}
	return exhibition
}

func (repo *ExhibitionRepository) Update(exhibition *domain.Exhibition, user int) *domain.Exhibition {
	params := make(map[string]string, 0)
	for _, v := range exhibition.Info {
		params[v.Type] = v.Value
	}
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdate, exhibition.Name, exhibition.Description, params, exhibition.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Exhibition updating error:", exhibition.ID, err, result.RowsAffected())
		return nil
	}
	return exhibition
}

func (repo *ExhibitionRepository) UpdateImage(exhibition *domain.Exhibition, user int) *domain.Exhibition {
	result, err := repo.db.Pool.Exec(context.Background(), queryUpdateImage,
		exhibition.Image, exhibition.Sizes.Height, exhibition.Sizes.Width, exhibition.ID, user)
	if err != nil || result.RowsAffected() == 0 {
		log.Println("Exhibition image updating error:", exhibition.ID, err, result.RowsAffected())
		return nil
	}
	return exhibition
}

func (repo *ExhibitionRepository) Delete(id, user int) error {
	result, err := repo.db.Pool.Exec(context.Background(), queryDeleteID, id, user)
	if err != nil {
		log.Println("Exhibition deleting error:", id, err)
		return err
	} else if result.RowsAffected() == 0 {
		log.Println("Exhibition for deleting not found")
		return errors.New("Exhibition not found")
	}
	return nil
}
