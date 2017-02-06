package main

const APPSecretKey  = "secret_key123456"

func main() {
	app := NewApp()

	app.ListenAndServe()
}