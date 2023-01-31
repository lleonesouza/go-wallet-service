package handlers

import (
	"net/http"
	"q2bank/handlers/dtos"
	"q2bank/services"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ShopkeeperHandler struct {
	service *services.Services
}

func (s *ShopkeeperHandler) Get(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)

	shopkeeper, err := s.service.Shopkeeper.Get(claims.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, shopkeeper)
}

func (s *ShopkeeperHandler) Create(c echo.Context) error {
	_shopkeeper := new(dtos.CreateShopkeeperDTO)
	if err := c.Bind(_shopkeeper); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	shopkeeper, err := s.service.Shopkeeper.Create(_shopkeeper)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, shopkeeper)
}

func (s *ShopkeeperHandler) Update(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)

	_shopkeeper := new(dtos.UpdateShopkeeperDTO)
	if err := c.Bind(_shopkeeper); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	shopkeeper, err := s.service.Shopkeeper.Update(claims.ID, _shopkeeper)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, shopkeeper)
}

func (s *ShopkeeperHandler) Login(c echo.Context) error {
	shopkeeper := new(dtos.LoginShopkeeperDTO)
	if err := c.Bind(shopkeeper); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	token, err := s.service.Shopkeeper.Login(shopkeeper.Email, shopkeeper.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
