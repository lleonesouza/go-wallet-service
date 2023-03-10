package handlers

import (
	"bff-answerfy/config"
	"bff-answerfy/errors"
	"bff-answerfy/handlers/dtos"
	"bff-answerfy/services"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type UserHandler struct {
	service   *services.Services
	env       *config.Envs
	validator *validator.Validate
	errors    *errors.Errors
}

//	@Description	Get account information.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.UserResponseDTO
//	@Failure		401	{object}	dtos.UnauthorizedError
//	@Failure		500	{object}	dtos.GeneralError
//	@Security		ApiKeyAuth
//	@Router			/user [get]
func (u *UserHandler) Get(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	user, err := u.service.User.GetById(claims.ID)
	if err != nil {
		unauthorizedError := u.errors.UnauthorizedError()
		c.JSON(unauthorizedError.Status, unauthorizedError)
	}

	return c.JSON(http.StatusOK, u.service.User.Filter(user))
}

//	@Description	Create a User account.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.CreateUserDTO	true	"Create User Account Input"
//	@Success		201		{object}	dtos.UserResponseDTO
//	@Failure		400		{object}	dtos.GeneralError
//	@Failure		500		{object}	dtos.GeneralError
//	@Router			/user [post]
func (u *UserHandler) Create(c echo.Context) error {
	// Bind User
	user := new(dtos.CreateUserDTO)
	err := c.Bind(user)
	if err != nil {
		typeError := u.errors.TypeError(err.Error())
		return c.JSON(typeError.Status, typeError)
	}

	// Validate body
	err = u.validator.Struct(user)
	if err != nil {
		bodyErr := u.errors.BodyError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	// Check if email exists
	_, err = u.service.User.GetByEmail(user.Email)
	if err == nil {
		conflictErr := u.errors.EmailRegistered(user.Email)
		return c.JSON(conflictErr.Status, conflictErr)
	}

	// Check if CPF exists
	_, err = u.service.User.GetByCPF(user.CPF)
	if err == nil {
		conflictErr := u.errors.CpfRegistered(user.CPF)
		return c.JSON(conflictErr.Status, conflictErr)
	}

	// Create Wallet
	wallet, err := u.service.Wallet.Create()
	if err != nil {
		creationErr := u.errors.CreateWalletError(err.Error())
		return c.JSON(creationErr.Status, creationErr)
	}

	// Create User
	response, err := u.service.User.Create(user, wallet)
	if err != nil {
		creationErr := u.errors.CreateUserError(err.Error())
		return c.JSON(creationErr.Status, creationErr)
	}

	return c.JSON(http.StatusCreated, u.service.User.Filter(response))
}

//	@Description	Update 'Name' and/or 'Lastname' of User account.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.UpdateUserDTO	true	"User"
//	@Success		200		{object}	dtos.UserResponseDTO
//	@Failure		401		{object}	dtos.UnauthorizedError
//	@Failure		400		{object}	dtos.GeneralError
//	@Security		ApiKeyAuth
//	@Router			/user [put]
func (u *UserHandler) Update(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	// Bind User
	_user := new(dtos.UpdateUserDTO)
	err := c.Bind(_user)
	if err != nil {
		typeErr := u.errors.TypeError(err.Error())
		return c.JSON(typeErr.Status, typeErr)
	}

	// Validate Body
	err = u.validator.Struct(_user)
	if err != nil {
		bodyErr := u.errors.BodyError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	// Update
	response, err := u.service.User.Update(claims.ID, _user)
	if err != nil {
		bodyErr := u.errors.UpdateUserError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	return c.JSON(http.StatusOK, u.service.User.Filter(response))
}

//	@Description	Login
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.LoginUserDTO	true	"User"
//	@Success		200		{object}	dtos.LoginResponseDTO
//	@Failure		400		{object}	dtos.GeneralError
//	@Router			/user/login [post]
func (u *UserHandler) Login(c echo.Context) error {
	// Bind LoginUserDTO
	user := new(dtos.LoginUserDTO)
	err := c.Bind(user)
	if err != nil {
		typeErr := u.errors.TypeError(err.Error())
		return c.JSON(typeErr.Status, typeErr)
	}

	// Validate LoginUserDTO
	err = u.validator.Struct(user)
	if err != nil {
		bodyErr := u.errors.BodyError(err.Error())
		return c.JSON(bodyErr.Status, bodyErr)
	}

	// Get User By Email
	completeUser, err := u.service.User.GetByEmail(user.Email)
	if err != nil {
		conflictErr := u.errors.EmailRegistered(user.Email)
		return c.JSON(conflictErr.Status, conflictErr)
	}

	// Compare Password
	err = u.service.User.CheckPasswordHash(user.Password, completeUser.Password)
	if err != nil {
		unauthorizedErr := u.errors.UnauthorizedError()
		return c.JSON(unauthorizedErr.Status, unauthorizedErr)
	}

	// Login
	token, err := u.service.User.Login(completeUser)
	if err != nil {
		loginErr := u.errors.LoginError()
		return c.JSON(loginErr.Status, loginErr)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
