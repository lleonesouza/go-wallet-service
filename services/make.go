package services

import (
	"context"
	"fmt"
	"q2bank/prisma/db"
)

type Services struct {
	User        *User
	Shopkeeper  *Shopkeeper
	Transaction *Transaction
}

func MakeServices() *Services {
	client := db.NewClient()
	err := client.Prisma.Connect()

	if err != nil {
		fmt.Errorf("%v", err.Error())
		panic(err)
	}

	// defer func() {
	// 	if err := client.Prisma.Disconnect(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	ctx := context.Background()

	wallet := &Wallet{client, ctx}

	return &Services{
		User: &User{
			client: client,
			wallet: wallet,
			ctx:    ctx,
		},
		Shopkeeper: &Shopkeeper{
			client: client,
			wallet: wallet,
			ctx:    ctx,
		},
		Transaction: &Transaction{
			client: client,
			ctx:    ctx,
		},
	}
}
