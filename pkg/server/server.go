package server

import (
	"ar_exhibition/pkg/database"
	"os"
	"strconv"

	"github.com/aerogo/aero"
)

func Run(Configure func(*aero.Application, interface{}) *aero.Application,
	funcs ...func(interface{}) interface{}) {
	app := aero.New()

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
