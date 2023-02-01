package services

import (
	"context"
	"q2bank/config"
	"q2bank/prisma/db"
)

type Services struct {
	User        *User
	Shopkeeper  *Shopkeeper
	Transaction *Transaction
}

func MakeServices(env *config.Envs) *Services {
	client := db.NewClient()
	err := client.Prisma.Connect()

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	wallet := &Wallet{client, ctx, env}

	return &Services{
		User: &User{
			client: client,
			wallet: wallet,
			ctx:    ctx,
			env:    env,
		},
		Shopkeeper: &Shopkeeper{
			client: client,
			wallet: wallet,
			ctx:    ctx,
			env:    env,
		},
		Transaction: &Transaction{
			client: client,
			ctx:    ctx,
			env:    env,
		},
	}
}
