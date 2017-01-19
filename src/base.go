package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"errors"
	"log"
)

type Context struct {
	response http.ResponseWriter
	request *http.Request
	query url.Values
}

type APIResult map[string]interface{}

type APIMethods map[string]func(*Context) APIResult

type APIModule interface {
	Name() string
	Export() APIMethods
	Dispose()
}

var apiFatalErr = errors.New("api fatal")

func (ctx *Context) Fatal(msg string, err error) {
	if err!=nil {
		log.Printf("%s %s %s %s", ctx.request.Method, ctx.request.RequestURI, msg, err.Error())
	}else{
		log.Printf("%s %s %s", ctx.request.Method, ctx.request.RequestURI,msg)
	}
	http.Error(ctx.response, `{"error":"}`+msg+`"}`, 500)
	panic(apiFatalErr)
}

func (ctx *Context) Get(name string) string {
	if ctx.query == nil {
		ctx.query = ctx.request.URL.Query()
	}
	return ctx.query[name][0]
}

func APIHandler(logic func(ctx *Context) APIResult) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil && err != apiFatalErr {
				log.Printf("%s PANIC!! - %v", r.RequestURI, err)
				http.Error(w, `{"error":"Internal error"}`, 500)
			}
		}()

		ctx := Context{response: w, request: r}

		result := logic(&ctx)

		data, err := json.Marshal(result)
		if err != nil {
			ctx.Fatal("JSON marshal failed", fmt.Errorf("%v, %s", r, err))
		}
		ctx.response.Write(data)
	})
}