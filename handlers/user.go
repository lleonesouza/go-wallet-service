package handlers

import (
	"net/http"
	"q2bank/config"
	"q2bank/handlers/dtos"
	"q2bank/services"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *services.Services
	env     *config.Envs
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

	user, err := u.service.User.Get(claims.ID)

	filteredUser := u.service.User.Filter(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, filteredUser)
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
	_user := new(dtos.CreateUserDTO)
	if err := c.Bind(_user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	user, err := u.service.User.Create(_user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	filteredUser := u.service.User.Filter(user)

	return c.JSON(http.StatusCreated, filteredUser)
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

	_user := new(dtos.UpdateUserDTO)
	if err := c.Bind(_user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	user, err := u.service.User.Update(claims.ID, _user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	m := make(map[string]string)
	m["name"] = user.Name
	m["lastname"] = user.Lastname
	m["id"] = user.ID
	m["create_at"] = user.CreatedAt.String()
	m["update_at"] = user.UpdatedAt.String()

	return c.JSON(http.StatusOK, m)
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
	user := new(dtos.LoginUserDTO)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	token, err := u.service.User.Login(user.Email, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
