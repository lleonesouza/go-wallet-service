package services

import (
	"context"
	"q2bank/prisma/db"
)

type Wallet struct {
	client *db.PrismaClient
	ctx    context.Context
}

func (w *Wallet) Create() (*db.WalletModel, error) {
	address := GenerateAddress()

	wallet, err := w.client.Wallet.CreateOne(
		db.Wallet.Address.Set(address),
		db.Wallet.Balance.Set(50),
	).Exec(w.ctx)

	return wallet, err
}
