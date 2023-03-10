package services

import (
	"bff-answerfy/config"
	"bff-answerfy/prisma/db"
	"context"
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

func (w *Wallet) AddCoins() (*db.WalletModel, error) {
	initialBalance, err := strconv.Atoi(w.env.WALLET_BALANCE_INIT)
	if err != nil {
		return nil, err
	}

	wallet, err := w.client.Wallet.CreateOne(
		db.Wallet.Balance.Set(initialBalance),
	).Exec(w.ctx)

	return wallet, err
}

func (w *Wallet) RemoveCoins() (*db.WalletModel, error) {
	initialBalance, err := strconv.Atoi(w.env.WALLET_BALANCE_INIT)
	if err != nil {
		return nil, err
	}

	wallet, err := w.client.Wallet.CreateOne(
		db.Wallet.Balance.Set(initialBalance),
	).Exec(w.ctx)

	return wallet, err
}

func (w *Wallet) GetWallet(user_id string) (*db.WalletModel, error) {
	user, err := w.client.User.FindUnique(
		db.User.ID.Equals(user_id),
	).With(
		db.User.Wallet.Fetch(),
	).Exec(w.ctx)

	if err != nil {
		return nil, err
	}

	return user.Wallet(), nil
}
