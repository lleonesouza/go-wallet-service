package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"q2bank/handlers/dtos"
	"q2bank/prisma/db"
)

type VerifyTransaction struct {
	Authorization bool `json:"authorization"`
}

type Transaction struct {
	client *db.PrismaClient
	ctx    context.Context
}

func (tx *Transaction) List(user_id string) ([]db.TransactionsModel, error) {
	user, err := tx.client.User.FindUnique(
		db.User.ID.Equals(user_id),
	).With(db.User.Wallet.Fetch().With(
		db.Wallet.Transactions.Fetch().With(db.Transactions.Wallets.Fetch()),
	)).Exec(tx.ctx)

	if err != nil {
		return nil, err
	}

	transactions := user.Wallet().Transactions()

	return transactions, err
}

func (tx *Transaction) getUserWallet(id string) (*db.WalletModel, error) {
	user, err := tx.client.User.FindUnique(
		db.User.ID.Equals(id),
	).With(
		db.User.Wallet.Fetch(),
	).Exec(tx.ctx)

	if err != nil {
		return nil, fmt.Errorf("Transaction not completed, User with %s not found", id)
	}

	wallet := user.Wallet()

	return wallet, nil
}

func (tx *Transaction) getUserWalletByEmail(email string) (*db.WalletModel, error) {
	user, err := tx.client.User.FindUnique(
		db.User.Email.Equals(email),
	).With(
		db.User.Wallet.Fetch(),
	).Exec(tx.ctx)

	if err != nil {
		return nil, fmt.Errorf("Transaction not completed, User with email '%s' not found", email)
	}

	wallet := user.Wallet()

	return wallet, nil
}

func (tx *Transaction) getShopkeeperWalletByEmail(email string) (*db.WalletModel, error) {
	shopkeeper, err := tx.client.Shopkeeper.FindUnique(
		db.Shopkeeper.ID.Equals(email),
	).With(
		db.Shopkeeper.Wallet.Fetch(),
	).Exec(tx.ctx)

	if err != nil {
		return nil, fmt.Errorf("Transaction not completed, Shopkeeper with Email %s not found", email)
	}

	wallet := shopkeeper.Wallet()

	return wallet, nil
}

func (tx *Transaction) getWallet(email string) (*db.WalletModel, error) {
	userWallet, _ := tx.getUserWalletByEmail(email)

	if userWallet != nil {
		return userWallet, nil
	}

	shopkeeperWallet, _ := tx.getShopkeeperWalletByEmail(email)

	if shopkeeperWallet != nil {
		return userWallet, nil
	}

	return nil, fmt.Errorf("Transaction not completed, receiver '%s' not found", email)
}

func (tx *Transaction) validTransaction() (*VerifyTransaction, error) {
	resp, err := http.Get("https://run.mocky.io/v3/d02168c6-d88d-4ff2-aac6-9e9eb3425e31")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := &VerifyTransaction{}
	json.Unmarshal(body, &data)

	return data, nil
}

func (tx *Transaction) isNormalUser(userType string) error {
	if userType != "user" {
		return errors.New("Only normal users can make transactions")
	}
	return nil
}

func (tx *Transaction) Filter(
	transaction *db.TransactionsModel,
	to string,
	from string,
) *dtos.ResponseTransactionDTO {
	return &dtos.ResponseTransactionDTO{
		Id:           transaction.ID,
		Value:        transaction.Value,
		FromWalletId: from,
		ToWalletId:   to,
		CreateAt:     transaction.CreatedAt.String(),
		UpdateAt:     transaction.CreatedAt.String(),
	}
}

func (tx *Transaction) CreateTransaction(
	fromUserId string,
	userType string,
	_tx dtos.CreateTransactionDTO,
) (*dtos.ResponseTransactionDTO, error) {
	err := tx.isNormalUser(userType)
	if err != nil {
		return nil, err
	}

	fromUserWallet, err := tx.getUserWallet(fromUserId)
	if err != nil {
		return nil, err
	}

	if _tx.Value > fromUserWallet.Balance {
		return nil, errors.New("Insufficient Coins")
	}

	toWallet, err := tx.getWallet(_tx.To)
	if err != nil {
		return nil, err
	}

	validTx, err := tx.validTransaction()
	if err != nil {
		return nil, err
	}

	if !validTx.Authorization {
		return nil, errors.New("Transaction is not authorized")
	}

	transaction := tx.client.Transactions.CreateOne(
		db.Transactions.Value.Set(_tx.Value),
		db.Transactions.Wallets.Link(
			db.Wallet.ID.Equals(fromUserWallet.ID),
			db.Wallet.ID.Equals(toWallet.ID),
		),
	).With(db.Transactions.Wallets.Fetch()).Tx()

	FromWallet := tx.client.Wallet.FindUnique(
		db.Wallet.ID.Equals(fromUserWallet.ID),
	).Update(
		db.Wallet.Balance.Decrement(_tx.Value),
	).Tx()

	ToWallet := tx.client.Wallet.FindUnique(
		db.Wallet.ID.Equals(toWallet.ID),
	).Update(
		db.Wallet.Balance.Increment(_tx.Value),
	).Tx()

	err = tx.client.Prisma.Transaction(transaction, FromWallet, ToWallet).Exec(tx.ctx)
	if err != nil {
		return nil, err
	}

	response := tx.Filter(transaction.Result(), toWallet.ID, fromUserWallet.ID)

	return response, nil
}
