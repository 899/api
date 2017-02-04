package main

import "time"

type HelperAPI struct {
}

func (app *App) newHelperAPI() *HelperAPI {
	return &HelperAPI{}
}

func (api *HelperAPI) Dispose() {
}

func (api *HelperAPI) Name() string {
	return "helper"
}

func (api *HelperAPI) Export() APIMethods {
	return APIMethods{
		"test": api.Test,
	}
}

func (api *HelperAPI) Test(ctx *Context) APIResult {
	return APIResult{
		"time": time.Now().UTC().Unix(),
		"result": "help api test",
	}
}
