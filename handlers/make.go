package handlers

import (
	"bff-answerfy/config"
	"bff-answerfy/errors"
	"bff-answerfy/handlers/validation"
	"bff-answerfy/services"
)

type Handlers struct {
	User     UserHandler
	Wallet   WalletHandler
	Question QuestionHandler
}

func MakeHandlers(env *config.Envs) *Handlers {

	validator := validation.MakeValidator()
	services := services.MakeServices(env)
	errors := errors.MakeErrors()

	return &Handlers{
		User: UserHandler{
			service:   services,
			env:       env,
			validator: validator,
			errors:    errors,
		},
		Wallet: WalletHandler{
			service:   services,
			env:       env,
			validator: validator,
			errors:    errors,
		},
		Question: QuestionHandler{
			service:   services,
			env:       env,
			validator: validator,
			errors:    errors,
		},
	}
}
