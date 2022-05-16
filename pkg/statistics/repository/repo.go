package repository

import (
	"ar_exhibition/pkg/database"
	"ar_exhibition/pkg/domain"
	"context"
	"fmt"
	"log"
)

const (
	querySave   = `insert into stats (port, method, status, urls, perform) values ($1, $2, $3, $4, $5);`
	queryGetAll = `select port, method, status, urls, perform, reqtime from stats;`
	queryGet    = `select port, method, status, urls, perform, reqtime from stats where `
)

type StatRepository struct {
	db *database.DBManager
}

func StatRepo(db interface{}) interface{} {
	instance, ok := db.(*database.DBManager)
	if ok {
		return &StatRepository{db: instance}
	}
	log.Println("Unknown object instead of db-manager")
	return nil
}

func (repo *StatRepository) Save(stat *domain.Stats) {
	_, err := repo.db.Pool.Exec(context.Background(), querySave, stat.Port, stat.Method, stat.Status, stat.URL, stat.Duration)
	if err != nil {
		log.Println("Error while saving stat to db:", err)
	}
}

func (repo *StatRepository) Get(port, status int, method string) []*domain.Stats {
	counter, sql := 1, queryGet
	params := make([]interface{}, 0)
	if port > 0 {
		params = append(params, port)
		sql += fmt.Sprintf("port = $%d", counter)
		counter++
	}
	if status != 0 {
		if counter > 1 {
			sql += ","
		}
		params = append(params, status)
		sql += fmt.Sprintf("status = $%d", counter)
		counter++
	}
	if method != "" {
		if counter > 1 {
			sql += ","
		}
		params = append(params, method)
		sql += fmt.Sprintf("method = $%d", counter)
	}
	sql += ";"

	rows, err := repo.db.Pool.Query(context.Background(), sql, params...)
	if err != nil {
		log.Println("Cannot get stats with params:", err)
		return nil
	}
	defer rows.Close()

	result := make([]*domain.Stats, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		stat := new(domain.Stats)
		if rows.Scan(&stat.Port, &stat.Method, &stat.Status, &stat.URL, &stat.Duration, &stat.When) != nil {
			return result
		}
		result = append(result, stat)
	}
	return result
}

func (repo *StatRepository) GetAll() []*domain.Stats {
	rows, err := repo.db.Pool.Query(context.Background(), queryGetAll)
	if err != nil {
		log.Println("Cannot get all stats:", err)
		return nil
	}
	defer rows.Close()

	result := make([]*domain.Stats, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		stat := new(domain.Stats)
		if rows.Scan(&stat.Port, &stat.Method, &stat.Status, &stat.URL, &stat.Duration, &stat.When) != nil {
			return result
		}
		result = append(result, stat)
	}
	return result
}
