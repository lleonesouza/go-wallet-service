package handlers

import (
	"q2bank/config"
	"q2bank/errors"
	"q2bank/handlers/validation"
	"q2bank/services"
)

type Handlers struct {
	Shopkeeper  ShopkeeperHandler
	User        UserHandler
	Transaction TransactionHandler
}

func MakeHandlers(env *config.Envs) *Handlers {

	validator := validation.MakeValidator()
	services := services.MakeServices(env)
	errors := errors.MakeErrors()

	return &Handlers{
		Shopkeeper: ShopkeeperHandler{
			service:   services,
			env:       env,
			validator: validator,
			errors:    errors,
		},
		User: UserHandler{
			service:   services,
			env:       env,
			validator: validator,
			errors:    errors,
		},
		Transaction: TransactionHandler{
			service:   services,
			env:       env,
			validator: validator,
			errors:    errors,
		},
	}
}
