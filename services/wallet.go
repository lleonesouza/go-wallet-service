package services

import (
	"context"
	"q2bank/config"
	"q2bank/prisma/db"
	"strconv"
)

type Wallet struct {
	client *db.PrismaClient
	ctx    context.Context
	env    *config.Envs
}

func (w *Wallet) Create() (*db.WalletModel, error) {
	initialBalance, err := strconv.Atoi(w.env.WALLET_BALANCE_INIT)
	if err != nil {
		return nil, err
	}

	wallet, err := w.client.Wallet.CreateOne(
		db.Wallet.Balance.Set(initialBalance),
	).Exec(w.ctx)

	return wallet, err
}
