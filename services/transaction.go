package services

import (
	"bff-answerfy/config"
	"bff-answerfy/prisma/db"
	"context"
)

type Transaction struct {
	client *db.PrismaClient
	ctx    context.Context
	env    *config.Envs
}

func (tx *Transaction) List(wallet_id string) ([]db.TransactionModel, error) {
	txs, err := tx.client.Transaction.FindMany(
		db.Transaction.WalletID.Equals(wallet_id),
	).Exec(tx.ctx)

	if err != nil {
		return nil, err
	}

	return txs, err
}

func (tx *Transaction) Create(txType string, wallet_id string, value int) (any, error) {
	if txType == "add" {
		transaction, err := tx.AddTx(wallet_id, value)
		if err != nil {
			return nil, err
		}
		return transaction, nil
	}

	hasCredit, err := tx.hasCredit(wallet_id, value)
	if err != nil {
		return nil, err
	}

	if !hasCredit {
		return nil, err
	}

	if txType == "remove" {
		transaction, err := tx.RmTx(wallet_id, value)
		if err != nil {
			return nil, err
		}
		return transaction, nil
	}

	return nil, err
}

func (tx *Transaction) AddTx(wallet_id string, value int) (*db.TransactionModel, error) {
	transaction := tx.client.Transaction.CreateOne(
		db.Transaction.Value.Set(value),
		db.Transaction.Type.Set("add"),
		db.Transaction.Wallet.Link(
			db.Wallet.ID.Equals(wallet_id),
		),
	).With(db.Transaction.Wallet.Fetch()).Tx()

	w := tx.client.Wallet.FindUnique(
		db.Wallet.ID.Equals(wallet_id),
	).Update(
		db.Wallet.Balance.Increment(value),
	).Tx()

	err := tx.client.Prisma.Transaction(transaction, w).Exec(tx.ctx)
	if err != nil {
		return nil, err
	}

	return transaction.Result(), nil
}

func (tx *Transaction) RmTx(wallet_id string, value int) (*db.TransactionModel, error) {
	transaction := tx.client.Transaction.CreateOne(
		db.Transaction.Value.Set(value),
		db.Transaction.Type.Set("remove"),
		db.Transaction.Wallet.Link(
			db.Wallet.ID.Equals(wallet_id),
		),
	).With(db.Transaction.Wallet.Fetch()).Tx()

	w := tx.client.Wallet.FindUnique(
		db.Wallet.ID.Equals(wallet_id),
	).Update(
		db.Wallet.Balance.Decrement(value),
	).Tx()

	err := tx.client.Prisma.Transaction(transaction, w).Exec(tx.ctx)
	if err != nil {
		return nil, err
	}

	return transaction.Result(), nil
}

func (tx *Transaction) hasCredit(wallet_id string, value int) (bool, error) {
	w, err := tx.client.Wallet.FindUnique(
		db.Wallet.ID.Equals(wallet_id),
	).Exec(tx.ctx)

	if err != nil {
		return false, err
	}

	if w.Balance > value {
		return true, nil
	}

	return false, nil
}
