package server

import (
	"ar_exhibition/pkg/database"

	"github.com/aerogo/aero"
)

func Run(Configure func(*aero.Application, interface{}) *aero.Application,
	funcs ...func(interface{}) interface{}) {
	app := aero.New()
	db := database.Connect()

	var repo, usecases, handlers interface{}
	switch len(funcs) {
	case 3:
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
