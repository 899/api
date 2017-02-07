package main

import (
	"time"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

type UserAPI struct {
}

func (app *App) newUserAPI() *UserAPI {
	return &UserAPI{}
}

func (api *UserAPI) Dispose() {
}

func (api *UserAPI) Name() string {
	return "user"
}

func (api *UserAPI) Export() RegisterMethods {
	return AppBootstrap(
		RegisterMethods{
			"login":RegisterMethod{api.login,Bootstrap{false}},
			"test":RegisterMethod{api.Test,Bootstrap{true}},
			"logout":RegisterMethod{api.logout,Bootstrap{true}},
			"index":RegisterMethod{api.index,Bootstrap{true}},
			"refresh":RegisterMethod{api.refreshToken,Bootstrap{true}},
		})
}

type Key int

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}


func (api *UserAPI) Test(ctx *Context) APIResult {
	return APIResult{
		"result": "user api test",
	}
	//obj := HelperAPI{}
	//res := obj.Test(ctx)
	//return  res
}

func (api *UserAPI) setToken(ctx *Context) {
	res := ctx.response
	//req := ctx.request
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	claims := MyClaims{
		"username",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:9000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString([]byte(APPSecretKey))

	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(res, &cookie)

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