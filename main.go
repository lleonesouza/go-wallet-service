package main

import (
	"bff-answerfy/config"
	"bff-answerfy/handlers"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	_ "bff-answerfy/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			bff-answerfy API
//	@version		1.0
//	@description	bff-answerfy Challenge

//	@contact.name	leone de souza
//	@contact.url	https://github.com/lleonesouza
//	@contact.email	lleonesouza@live.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:1323
//	@BasePath	/

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Token used authenticate 'User' and 'Shopkeeper'

func main() {
	e := echo.New()
	envs := config.MakeEnvs()
	h := handlers.MakeHandlers(envs)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(config.JwtCustomClaims)
		},
		SigningKey: []byte(envs.JWT_SECRET),
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	user := e.Group("/user")
	user.POST("", h.User.Create)
	user.POST("/login", h.User.Login)
	user.GET("", h.User.Get, echojwt.WithConfig(config))
	user.PUT("", h.User.Update, echojwt.WithConfig(config))

	question := e.Group("/question")
	question.POST("", h.Question.Create, echojwt.WithConfig(config))
	question.GET("", h.Question.List, echojwt.WithConfig(config))

	wallet := e.Group("/wallet")
	wallet.GET("", h.Wallet.Get, echojwt.WithConfig(config))

	e.Logger.Fatal(e.Start(":1323"))
}
