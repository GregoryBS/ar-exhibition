package server

import (
	"ar_exhibition/pkg/domain"
	"ar_exhibition/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aerogo/aero"
	"github.com/go-redis/redis"
)

func PanicRecover(next aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic with err: '%s' on url: %s\n", err, ctx.Path())
				ctx.Error(http.StatusInternalServerError)
			}
		}()
		return next(ctx)
	}
}

func Logging(next aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		start := time.Now()
		err := next(ctx)
		end := time.Since(start)
		log.Println(ctx.Request().Method(), ctx.Request().Internal().RequestURI, ctx.Status(), end)

		stats := &domain.Stats{
			Port:     port,
			Method:   ctx.Request().Method(),
			Status:   ctx.Status(),
			URL:      ctx.Request().Internal().RequestURI,
			Duration: int(end.Milliseconds()),
		}
		sendStats(stats)
		return err
	}
}

func sendStats(stats *domain.Stats) {
	msg, _ := json.Marshal(stats)
	err := client.RPush(utils.MsgTag, msg).Err()
	if err != nil {
		log.Println("Unable to send statistics in redis")
	}
}

var port, _ = strconv.Atoi(os.Getenv("PORT"))
var client = redis.NewClient(&redis.Options{
	Addr:     utils.RedisService,
	Password: "",
	DB:       0,
})
