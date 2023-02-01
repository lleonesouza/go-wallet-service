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

	return c.JSON(http.StatusOK, filteredUser)
}

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
