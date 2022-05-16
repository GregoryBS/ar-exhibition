package usecase

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/statistics/repository"
	"log"
)

type StatUsecase struct {
	repo *repository.StatRepository
}

func StatUsecases(repo interface{}) interface{} {
	instance, ok := repo.(*repository.StatRepository)
	if ok {
		return &StatUsecase{repo: instance}
	}
	log.Println("Unknown object instead of stat usecase")
	return nil
}

func (u *StatUsecase) SaveStat(stat *domain.Stats) {
	u.repo.Save(stat)
}

func (u *StatUsecase) GetStats(port, status int, method string) []*domain.Stats {
	if port > 0 || status != 0 || method != "" {
		return u.repo.Get(port, status, method)
	}
	return u.repo.GetAll()
}
