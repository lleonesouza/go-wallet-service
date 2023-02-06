package handlers

import (
	"net/http"
	"q2bank/config"
	"q2bank/errors"
	"q2bank/handlers/dtos"
	"q2bank/services"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type ShopkeeperHandler struct {
	service   *services.Services
	env       *config.Envs
	validator *validator.Validate
	errors    *errors.Errors
}

//	@Description	Get account information.
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.ShopkeeperResponseDTO
//	@Failure		401	{object}	dtos.UnauthorizedError
//	@Failure		500	{object}	dtos.GeneralError
//	@Security		ApiKeyAuth
//	@Router			/shopkeeper [get]
func (s *ShopkeeperHandler) Get(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	shopkeeper, err := s.service.Shopkeeper.GetById(claims.ID)
	if err != nil {
		unauthorizedError := s.errors.UnauthorizedError()
		c.JSON(unauthorizedError.Status, unauthorizedError)
	}

	return c.JSON(http.StatusOK, s.service.Shopkeeper.Filter(shopkeeper))
}

//	@Description	Create a Shopkeeper account.
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.CreateShopkeeperDTO	true	"Create Shopkeeper Account Input"
//	@Success		201		{object}	dtos.ShopkeeperResponseDTO
//	@Failure		400		{object}	dtos.GeneralError
//	@Failure		500		{object}	dtos.GeneralError
//	@Router			/shopkeeper [post]
func (s *ShopkeeperHandler) Create(c echo.Context) error {
	// Bind Shopkeeper
	shopkeeper := new(dtos.CreateShopkeeperDTO)
	err := c.Bind(shopkeeper)
	if err != nil {
		typeError := s.errors.TypeError(err.Error())
		return c.JSON(typeError.Status, typeError)
	}

	// Validate body
	err = s.validator.Struct(shopkeeper)
	if err != nil {
		bodyErr := s.errors.BodyError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	// Check if email exists
	_, err = s.service.Shopkeeper.GetByEmail(shopkeeper.Email)
	if err == nil {
		conflictErr := s.errors.EmailRegistered(shopkeeper.Email)
		return c.JSON(conflictErr.Status, conflictErr)
	}

	// Check if CNPJ exists
	_, err = s.service.User.GetByCPF(shopkeeper.CNPJ)
	if err == nil {
		conflictErr := s.errors.CpfRegistered(shopkeeper.CNPJ)
		return c.JSON(conflictErr.Status, conflictErr)
	}

	// Create Wallet
	wallet, err := s.service.Wallet.Create()
	if err != nil {
		creationErr := s.errors.CreateWalletError(err.Error())
		return c.JSON(creationErr.Status, creationErr)
	}

	// Create Shopkeeper
	response, err := s.service.Shopkeeper.Create(shopkeeper, wallet)
	if err != nil {
		creationErr := s.errors.CreateUserError(err.Error())
		return c.JSON(creationErr.Status, creationErr)
	}

	return c.JSON(http.StatusCreated, s.service.Shopkeeper.Filter(response))
}

//	@Description	Update 'Name' and/or 'Lastname' of Shopkeeper account.
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.UpdateShopkeeperDTO	true	"Shopkeeper"
//	@Success		200		{object}	dtos.ShopkeeperResponseDTO
//	@Failure		401		{object}	dtos.UnauthorizedError
//	@Failure		400		{object}	dtos.GeneralError
//	@Security		ApiKeyAuth
//	@Router			/shopkeeper [put]
func (s *ShopkeeperHandler) Update(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	// Bind User
	shopkeeper := new(dtos.UpdateUserDTO)
	err := c.Bind(shopkeeper)
	if err != nil {
		typeErr := s.errors.TypeError(err.Error())
		return c.JSON(typeErr.Status, typeErr)
	}

	// Validate Body
	err = s.validator.Struct(shopkeeper)
	if err != nil {
		bodyErr := s.errors.BodyError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	// Update
	response, err := s.service.User.Update(claims.ID, shopkeeper)
	if err != nil {
		bodyErr := s.errors.UpdateUserError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	return c.JSON(http.StatusOK, s.service.User.Filter(response))
}

//	@Description	Login
//	@Tags			shopkeeper
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.LoginShopkeeperDTO	true	"Shopkeeper"
//	@Success		200		{object}	dtos.LoginResponseDTO
//	@Failure		400		{object}	dtos.GeneralError
//	@Router			/shopkeeper/login [post]
func (s *ShopkeeperHandler) Login(c echo.Context) error {
	// Bind LoginShopeekerDTO
	shopkeeper := new(dtos.LoginShopkeeperDTO)
	err := c.Bind(shopkeeper)
	if err != nil {
		typeErr := s.errors.TypeError(err.Error())
		return c.JSON(typeErr.Status, typeErr)
	}

	// Validate Body
	err = s.validator.Struct(shopkeeper)
	if err != nil {
		bodyErr := s.errors.BodyError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	// Get Shopkeeper By Email
	completeShokeeper, err := s.service.Shopkeeper.GetByEmail(shopkeeper.Email)
	if err == nil {
		conflictErr := s.errors.EmailRegistered(shopkeeper.Email)
		return c.JSON(conflictErr.Status, conflictErr)
	}

	// Compare Password
	err = s.service.Shopkeeper.CheckPasswordHash(completeShokeeper.Password, shopkeeper.Password)
	if err != nil {
		unauthorizedErr := s.errors.UnauthorizedError()
		return c.JSON(unauthorizedErr.Status, unauthorizedErr)
	}

	// Login
	token, err := s.service.User.Login(completeShokeeper)
	if err != nil {
		loginErr := s.errors.LoginError()
		return c.JSON(loginErr.Status, loginErr)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
