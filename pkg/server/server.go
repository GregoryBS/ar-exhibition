package server

import (
	"ar_exhibition/pkg/database"

	"github.com/aerogo/aero"
)

func Run(Configure func(*aero.Application, interface{}) *aero.Application,
	stat func(interface{}),
	funcs ...func(interface{}) interface{}) {
	app := aero.New()
	app.Use(PanicRecover, Logging)

	var db *database.DBManager
	var repo, usecases, handlers interface{}
	switch len(funcs) {
	case 3:
		app.Config.Ports.HTTP = port
		db = database.Connect()
		repo = funcs[2](db)
		fallthrough
	case 2:
		usecases = funcs[1](repo)
		handlers = funcs[0](usecases)
	}
	if stat != nil {
		go stat(handlers)
	}

	app.OnEnd(func() {
		database.Disconnect(db)
	})
	app = Configure(app, handlers)
	app.Run()
}
