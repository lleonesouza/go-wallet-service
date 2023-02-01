package handlers

import (
	"net/http"
	"q2bank/config"
	"q2bank/handlers/dtos"
	"q2bank/services"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ShopkeeperHandler struct {
	service *services.Services
}

//	@Description	Get account information.
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.ShopkeeperResponseDTO
//	@Failure		401	{object}	error.Unauthorized
//	@Failure		500	{object}	error.General
//	@Security		ApiKeyAuth
//	@Router			/shopkeeper [get]
func (s *ShopkeeperHandler) Get(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	shopkeeper, err := s.service.Shopkeeper.Get(claims.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	filteredShopkeeper := s.service.Shopkeeper.Filter(shopkeeper)

	return c.JSON(http.StatusOK, filteredShopkeeper)
}

//	@Description	Create a Shopkeeper account.
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.CreateShopkeeperDTO	true	"Create Shopkeeper Account Input"
//	@Success		201		{object}	dtos.ShopkeeperResponseDTO
//	@Failure		400		{object}	error.General
//	@Failure		500		{object}	error.General
//	@Router			/shopkeeper [post]
func (s *ShopkeeperHandler) Create(c echo.Context) error {
	_shopkeeper := new(dtos.CreateShopkeeperDTO)
	if err := c.Bind(_shopkeeper); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	shopkeeper, err := s.service.Shopkeeper.Create(_shopkeeper)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	filteredShopkeeper := s.service.Shopkeeper.Filter(shopkeeper)

	return c.JSON(http.StatusOK, filteredShopkeeper)
}

//	@Description	Update 'Name' and/or 'Lastname' of Shopkeeper account.
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.UpdateShopkeeperDTO	true	"Shopkeeper"
//	@Success		200		{object}	dtos.ShopkeeperResponseDTO
//	@Failure		401		{object}	error.Unauthorized
//	@Failure		400		{object}	error.General
//	@Security		ApiKeyAuth
//	@Router			/shopkeeper [put]
func (s *ShopkeeperHandler) Update(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	_shopkeeper := new(dtos.UpdateShopkeeperDTO)
	if err := c.Bind(_shopkeeper); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	shopkeeper, err := s.service.Shopkeeper.Update(claims.ID, _shopkeeper)

	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	m := make(map[string]string)
	m["name"] = shopkeeper.Name
	m["lastname"] = shopkeeper.Lastname
	m["id"] = shopkeeper.ID
	m["create_at"] = shopkeeper.CreatedAt.String()
	m["update_at"] = shopkeeper.UpdatedAt.String()

	return c.JSON(http.StatusOK, m)
}

//	@Description	Login
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.LoginShopkeeperDTO	true	"Shopkeeper"
//	@Success		200		{object}	dtos.ShopkeeperResponseDTO
//	@Failure		400		{object}	error.General
//	@Router			/shopkeeper/login [post]
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
