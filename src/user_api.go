package main

import (
	"time"
	"fmt"
	"context"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

type UserAPI struct {
}

func (app *App) newUserAPI() *UserAPI {
	return &UserAPI{}
}

type Key int

const MyKey Key = 0

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (api *UserAPI) Dispose() {
}

func (api *UserAPI) Name() string {
	return "user"
}

func (api *UserAPI) Export() APIMethods {
	return APIMethods{
		"login": api.login,
		"logout": api.logout,
		"test": api.test,
		"index": api.index,
		"refresh": api.refreshToken,
	}
}


func (api *UserAPI) test(ctx *Context) APIResult {
	//return APIResult{
	//	"result": "user api test",
	//}
	obj := HelperAPI{}
	res := obj.Test(ctx)
	return  res
}

func (api *UserAPI) setToken(ctx *Context) {
	res := ctx.response
	//req := ctx.request
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	claims := Claims{
		"username",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:9000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString([]byte("secret_key123456"))

	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(res, &cookie)

}

func (api *UserAPI) validateToken(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			http.NotFound(res, req)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(res, req)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), MyKey, *claims)
			page(res, req.WithContext(ctx))
		} else {
			http.NotFound(res, req)
			return
		}
	})
}

func (api *UserAPI) login(ctx *Context) APIResult{
	api.setToken(ctx)
	//http.Redirect(ctx.response, ctx.request, "index", 307)
	return APIResult{
		"result":true,
	}
}

func (api *UserAPI) refreshToken(ctx *Context) APIResult {
	api.setToken(ctx)
	return APIResult{
		"result":true,
	}
}

func (api *UserAPI) index(ctx *Context) APIResult {
	//res :=ctx.response
	//req :=ctx.request
	//claims, ok := req.Context().Value(MyKey).(Claims)
	//if !ok {
	//	http.NotFound(res, req)
	//	return APIResult{}
	//}

	//fmt.Fprintf(res, "Hello %s", claims.Username)

	return APIResult{
		"result": true,
		"info": "hello index",
	}
}

func (api *UserAPI) logout(ctx *Context) APIResult {
	res := ctx.response
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(res, &deleteCookie)
	return APIResult{
		"result":true,
	}
}