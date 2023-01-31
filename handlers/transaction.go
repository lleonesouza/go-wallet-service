package handlers

import (
	"net/http"
	"q2bank/handlers/dtos"
	"q2bank/services"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	service *services.Services
}

func (tx *TransactionHandler) Create(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)

	txInput := new(dtos.CreateTransactionDTO)
	if err := c.Bind(txInput); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}

	txOutput, err := tx.service.Transaction.CreateTransaction(claims.ID, claims.Type, *txInput)

	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, txOutput)
}

func (tx *TransactionHandler) List(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)

	txOutputs, err := tx.service.Transaction.List(claims.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, txOutputs)
}
