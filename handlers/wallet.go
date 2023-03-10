package handlers

import (
	"bff-answerfy/config"
	"bff-answerfy/errors"
	"bff-answerfy/services"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type WalletHandler struct {
	service   *services.Services
	env       *config.Envs
	validator *validator.Validate
	errors    *errors.Errors
}

//	@Description	Get transaction from Wallet.
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.ResponseTransactionDTO
//	@Failure		401	{object}	dtos.UnauthorizedError
//	@Security		ApiKeyAuth
//	@Router			/wallet [get]
func (w *WalletHandler) Get(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	wallet, err := w.service.Wallet.GetWallet(claims.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, wallet)
}
