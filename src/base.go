package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"errors"
	"log"
	"fmt"
	"strings"
	jwt "github.com/dgrijalva/jwt-go"
)

type Context struct {
	response http.ResponseWriter
	request *http.Request
	query url.Values
}

type APIResult map[string]interface{}

type APIModule interface {
	Name() string
	Export() RegisterMethods
	Dispose()
}

var apiFatalErr = errors.New("api fatal")

func (ctx *Context) Fatal(msg string, err error) {
	if err!=nil {
		log.Printf("%s %s %s %s", ctx.request.Method, ctx.request.RequestURI, msg, err.Error())
	}else{
		log.Printf("%s %s %s", ctx.request.Method, ctx.request.RequestURI,msg)
	}
	http.Error(ctx.response, `{"error":"`+msg+`}"}`, 500)
	panic(apiFatalErr)
}

func (ctx *Context) Get(name string) string {
	if ctx.query == nil {
		ctx.query = ctx.request.URL.Query()
	}
	return ctx.query[name][0]
}

func APIHandler(logic func(ctx *Context) APIResult,bootstrap Bootstrap) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil && err != apiFatalErr {
				log.Printf("%s PANIC!! - %v", req.RequestURI, err)
				http.Error(res, `{"error":"Internal error"}`, 500)
			}
		}()
		ctx := Context{response: res, request: req}

		// logic
		result := CtxBootstrap(logic,ctx, bootstrap)
		data, err := json.Marshal(result)
		if err != nil {
			ctx.Fatal("JSON marshal failed", fmt.Errorf("%v, %s", req, err))
		}

		ctx.response.Write(data)
	})
}

type RegisterMethod struct {
	method func(*Context) APIResult
	bootstrap Bootstrap
}
type Bootstrap struct{
	isValid bool
}

type RegisterMethods map[string]RegisterMethod

func AppBootstrap(registers RegisterMethods) RegisterMethods {
	// register config
	return registers
}

func CtxBootstrap(logic func(ctx *Context) APIResult,ctx Context,bootstrap Bootstrap) APIResult{
	if bootstrap.isValid {
		JwtTokenValid(ctx.response,ctx.request)
	}
	result := logic(&ctx)
	return result
}
func JwtTokenValid(res http.ResponseWriter,req *http.Request){
	ctx := Context{response: res, request: req}
	// validate token
	cookie, err := req.Cookie("Auth")
	if err != nil {
		ctx.Fatal("auth cookie failed", fmt.Errorf("%v, %s", req, err))
	}
	splitCookie := strings.Split(cookie.String(), "Auth=")
	token, err := jwt.ParseWithClaims(splitCookie[1], &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(APPSecretKey), nil
	})
	if err != nil {
		ctx.Fatal("token parse failed", fmt.Errorf("%v, %s", req, err))
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Sprintf("%s",claims)
	} else {
		ctx.Fatal("token validate failed", fmt.Errorf("%v, %s", req, err))
	}
}
