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

func (api *UserAPI) Name() string {
	return "user"
}

func (api *UserAPI) Dispose() {
}

func (api *UserAPI) Export() APIMethods {
	return APIMethods{
		"settoken": api.setToken,
		"logout": api.logout,
	}
}

func (api *UserAPI) setToken(ctx *Context) APIResult {
	res := ctx.response
	//req := ctx.request
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	claims := Claims{
		"myusername",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:9000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := token.SignedString([]byte("secret"))

	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(res, &cookie)

	//http.Redirect(res, req, "/profile", 307)
	return APIResult{
		"result": true,
	}
}

func (api *UserAPI) validate(page http.HandlerFunc) http.HandlerFunc {
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

func (api *UserAPI) protectedProfile(res http.ResponseWriter, req *http.Request) {
	claims, ok := req.Context().Value(MyKey).(Claims)
	if !ok {
		http.NotFound(res, req)
		return
	}

	fmt.Fprintf(res, "Hello %s", claims.Username)
}

func (api *UserAPI) logout(ctx *Context) APIResult {
	res := ctx.response
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(res, &deleteCookie)
	return APIResult{
		"result":true,
	}
}