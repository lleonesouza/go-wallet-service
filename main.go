package main

import (
	"os"
	"q2bank/handlers"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	h := handlers.MakeHandlers()

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SECRET")),
	}

	user := e.Group("/user")
	user.POST("", h.User.Create)
	user.POST("/login", h.User.Login)
	user.GET("", h.User.Get, echojwt.WithConfig(config))
	user.PUT("", h.User.Update, echojwt.WithConfig(config))

	shopkeeper := e.Group("/shopkeeper")
	shopkeeper.POST("", h.Shopkeeper.Create)
	shopkeeper.POST("/login", h.Shopkeeper.Login)
	shopkeeper.GET("", h.Shopkeeper.Get, echojwt.WithConfig(config))
	shopkeeper.PUT("", h.Shopkeeper.Update, echojwt.WithConfig(config))

	tx := e.Group("/transaction")
	tx.POST("", h.Transaction.Create, echojwt.WithConfig(config))
	tx.GET("", h.Transaction.List, echojwt.WithConfig(config))

	e.Logger.Fatal(e.Start(":1323"))
}
