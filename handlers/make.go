package handlers

import "q2bank/services"

type Handlers struct {
	Shopkeeper  ShopkeeperHandler
	User        UserHandler
	Transaction TransactionHandler
}

func MakeHandlers() *Handlers {

	services := services.MakeServices()

	return &Handlers{
		Shopkeeper: ShopkeeperHandler{
			service: services,
		},
		User: UserHandler{
			service: services,
		},
		Transaction: TransactionHandler{
			service: services,
		},
	}
}
