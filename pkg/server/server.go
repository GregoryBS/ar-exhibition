package server

import (
	"ar_exhibition/pkg/database"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aerogo/aero"
)

func PanicRecover(next aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Recovered from panic with err: '%s' on url: %s\n", err, ctx.Path())
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
		fmt.Println(ctx.Request().Method(), ctx.Request().Internal().RequestURI, ctx.Status(), time.Since(start))
		return err
	}
}

func Run(Configure func(*aero.Application, interface{}) *aero.Application,
	funcs ...func(interface{}) interface{}) {
	app := aero.New()
	app.Use(PanicRecover, Logging)

	var db *database.DBManager
	var repo, usecases, handlers interface{}
	switch len(funcs) {
	case 3:
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		app.Config.Ports.HTTP = port
		db = database.Connect()
		repo = funcs[2](db)
		fallthrough
	case 2:
		usecases = funcs[1](repo)
		handlers = funcs[0](usecases)
	}

	app.OnEnd(func() {
		database.Disconnect(db)
	})
	app = Configure(app, handlers)
	app.Run()
}
