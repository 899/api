package main

import (
	"net/http"
	"fmt"
)

const API_VERSION = "1.0"

type App struct {
	modules []APIModule
}

func NewApp() *App {
	app := &App{}
	// modules
	app.modules = []APIModule{
		//app.newHelperAPI(),
		app.newUserAPI(),
	}

	return app
}

func (app *App) ListenAndServe(){

	// setting
	addr := ":9000"

	mux := http.NewServeMux()

	// route bind
	for _,module := range app.modules {
		registers := module.Export()
		for path, register := range registers {
			path = fmt.Sprintf("/%s/%s/%s", API_VERSION, module.Name(),path)
			mux.Handle(path,APIHandler(register.method,register.bootstrap))
		}
	}

	http.ListenAndServe(addr,mux)

	app.dispose()
}

func (app *App) dispose() {
	for _,module := range app.modules {
		module.Dispose()
	}

}