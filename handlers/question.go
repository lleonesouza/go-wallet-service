package handlers

import (
	"bff-answerfy/config"
	"bff-answerfy/errors"
	"bff-answerfy/handlers/dtos"
	"bff-answerfy/services"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type QuestionHandler struct {
	service   *services.Services
	env       *config.Envs
	validator *validator.Validate
	errors    *errors.Errors
}

//	@Description	Create a Question.
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.CreateTransactionDTO	true	"Create Transaction Input"
//	@Success		201		{object}	dtos.ResponseTransactionDTO
//	@Failure		400		{object}	dtos.GeneralError
//	@Security		ApiKeyAuth
//	@Router			/question [post]
func (tx *QuestionHandler) Create(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	body := new(dtos.CreateQuestionDTO)
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, services.FormatError(err))
	}
	// Get User
	user, err := tx.service.User.GetById(claims.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	// Create Question
	question, err := tx.service.Question.Create(user.ID, body.Text)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	// Send to AI
	res, err := tx.service.AI.GenerateAnswer(body.Text)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}
	if res.Choices[0].Text == "" {
		// Update Question
		question, err = tx.service.Question.Update(question.ID, res.Choices[0].Text, "failed")
		if err != nil {
			return c.JSON(http.StatusBadRequest, services.FormatError(err))
		}
		return c.JSON(http.StatusBadRequest, services.FormatError(fmt.Errorf("Generate Answer not completed")))
	}

	// Create Tx
	_, err = tx.service.Transaction.Create("remove", user.WalletID, 1)
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	// Update Question
	question, err = tx.service.Question.Update(question.ID, res.Choices[0].Text, "completed")
	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusCreated, question)
}

//	@Description	Get transaction from Wallet.
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.ResponseTransactionDTO
//	@Failure		401	{object}	dtos.UnauthorizedError
//	@Security		ApiKeyAuth
//	@Router			/questions [get]
func (tx *QuestionHandler) List(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*config.JwtCustomClaims)

	questions, err := tx.service.Question.GetAll(claims.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, services.FormatError(err))
	}

	return c.JSON(http.StatusOK, questions)
}
