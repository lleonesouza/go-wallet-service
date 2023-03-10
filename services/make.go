package services

import (
	"bff-answerfy/config"
	"bff-answerfy/prisma/db"
	"context"
)

type Services struct {
	User        *User
	Transaction *Transaction
	Wallet      *Wallet
	AI          *AI
	Question    *Question
}

func MakeServices(env *config.Envs) *Services {
	client := db.NewClient()
	err := client.Prisma.Connect()

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	http := &Http{env: env}
	wallet := &Wallet{client, ctx, env}
	ai := &AI{
		env:    env,
		client: client,
		http:   http,
		ctx:    ctx,
	}

	return &Services{
		User: &User{
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
		Wallet: wallet,
		AI:     ai,
		Question: &Question{
			client: client,
			wallet: wallet,
			ctx:    ctx,
			env:    env,
		},
	}
}
