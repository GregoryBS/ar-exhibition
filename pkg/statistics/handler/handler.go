package handler

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/statistics/usecase"
	"ar_exhibition/pkg/utils"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/aerogo/aero"
	"github.com/go-redis/redis"
)

type StatHandler struct {
	u *usecase.StatUsecase
}

func StatHandlers(usecases interface{}) interface{} {
	instance, ok := usecases.(*usecase.StatUsecase)
	if ok {
		return &StatHandler{u: instance}
	}
	log.Println("Unknown object instead of stat handler")
	return nil
}

func ConfigureStat(app *aero.Application, handlers interface{}) *aero.Application {
	h, ok := handlers.(*StatHandler)
	if ok {
		app.Get(utils.BaseStatApi, h.GetStats)
	}
	return app
}

func (h *StatHandler) GetStats(ctx aero.Context) error {
	url := ctx.Request().Internal().URL.Query()
	port, _ := strconv.Atoi(url.Get("port"))
	method := url.Get("method")
	status, _ := strconv.Atoi(url.Get("status"))
	stats := h.u.GetStats(port, status, method)
	return ctx.JSON(stats)
}

func ReceiveStats(handlers interface{}) {
	h, ok := handlers.(*StatHandler)
	if !ok {
		log.Println("Unknown object instead of stat handler")
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     utils.RedisService,
		Password: "",
		DB:       0,
	})

	for {
		if result, err := client.BLPop(0*time.Second, utils.MsgTag).Result(); err == nil {
			stats := new(domain.Stats)
			if err = json.Unmarshal([]byte(result[1]), stats); err == nil {
				h.u.SaveStat(stats)
			} else {
				log.Println("Invalid stats json", err)
			}
		} else {
			log.Println("Unable to get message from queue", err)
		}
	}
}
