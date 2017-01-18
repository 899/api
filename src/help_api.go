package main

import "time"

type HelperAPI struct {
}

func (app *App) newHelperAPI() *HelperAPI {
	return &HelperAPI{}
}

func (api *HelperAPI) Name() string {
	return "helper"
}

func (api *HelperAPI) Dispose() {
}

func (api *HelperAPI) Export() APIMethods {
	return APIMethods{
		"now": api.now,
	}
}

func (api *HelperAPI) now(ctx *Context) APIResult {
	return APIResult{
		"time": time.Now().UTC().Unix(),
	}
}
