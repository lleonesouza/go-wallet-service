package services

import (
	"context"
	"errors"
	"os"
	"q2bank/config"
	"q2bank/handlers/dtos"
	"q2bank/prisma/db"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/go-playground/validator.v9"
)

type Shopkeeper struct {
	client *db.PrismaClient
	wallet *Wallet
	ctx    context.Context
}

func (s *Shopkeeper) Filter(shopkeeper *db.ShopkeeperModel) *dtos.ShopkeeperResponseDTO {
	return &dtos.ShopkeeperResponseDTO{
		ID:       shopkeeper.ID,
		Balance:  shopkeeper.Wallet().Balance,
		Name:     shopkeeper.Name,
		Lastname: shopkeeper.Lastname,
		CNPJ:     shopkeeper.Cnpj,
		Email:    shopkeeper.Email,
		CreateAt: shopkeeper.CreatedAt.String(),
		UpdateAt: shopkeeper.UpdatedAt.String(),
	}
}

func (s *Shopkeeper) Create(_shopkeeper *dtos.CreateShopkeeperDTO) (*db.ShopkeeperModel, error) {
	err := s.validate(_shopkeeper)
	if err != nil {
		return nil, err
	}

	passwordHash, err := HashPassword(_shopkeeper.Password)
	if err != nil {
		return nil, err
	}

	wallet, err := s.wallet.Create()
	if err != nil {
		return nil, err
	}

	shopkeeper, err := s.client.Shopkeeper.CreateOne(
		db.Shopkeeper.Email.Set(_shopkeeper.Email),
		db.Shopkeeper.Cnpj.Set(_shopkeeper.CNPJ),
		db.Shopkeeper.Password.Set(passwordHash),
		db.Shopkeeper.Name.Set(_shopkeeper.Name),
		db.Shopkeeper.Lastname.Set(_shopkeeper.Lastname),
		db.Shopkeeper.Wallet.Link(
			db.Wallet.ID.Equals(wallet.ID),
		),
	).With(db.Shopkeeper.Wallet.Fetch()).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	return shopkeeper, nil
}

func (s *Shopkeeper) Update(id string, _shopkeeper *dtos.UpdateShopkeeperDTO) (*db.ShopkeeperModel, error) {
	shopkeeper, err := s.client.Shopkeeper.FindUnique(
		db.Shopkeeper.ID.Equals(id),
	).Update(
		db.Shopkeeper.Name.Set(_shopkeeper.Name),
		db.Shopkeeper.Lastname.Set(_shopkeeper.Lastname),
	).Exec(s.ctx)

	if err != nil {
		return shopkeeper, err
	}

	return shopkeeper, nil
}

func (s *Shopkeeper) Get(id string) (*db.ShopkeeperModel, error) {
	shopkeeper, err := s.client.Shopkeeper.FindUnique(
		db.Shopkeeper.ID.Equals(id),
	).With(
		db.Shopkeeper.Wallet.Fetch(),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	return shopkeeper, nil
}

func (s *Shopkeeper) Login(email string, password string) (string, error) {
	shopkeeper, err := s.getByEmail(email)
	if err != nil {
		return "", errors.New("Email or Password incorrect")
	}

	if !CheckPasswordHash(password, shopkeeper.Password) {
		return "", errors.New("Email or Password incorrect")
	}

	claims := &config.JwtCustomClaims{
		Email: shopkeeper.Email,
		ID:    shopkeeper.ID,
		Type:  "shopkeeper",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))

	return t, nil
}

func (s *Shopkeeper) validate(shopkeeper *dtos.CreateShopkeeperDTO) error {
	v := validator.New()
	err := v.Struct(shopkeeper)

	if err != nil {
		return err
	}

	err = s.emailExists(shopkeeper.Email)
	if err != nil {
		return err
	}
	err = s.cnpjExists(shopkeeper.CNPJ)
	if err != nil {
		return err
	}

	return nil
}

func (s *Shopkeeper) emailExists(email string) error {
	_, err := s.client.Shopkeeper.FindUnique(db.Shopkeeper.Email.Equals(email)).Exec(s.ctx)
	if err == nil {
		return errors.New("The field 'Email' is already in use. Please use another email.")
	}
	return nil
}

func (s *Shopkeeper) validateShopkeeper(shopkeeper *dtos.CreateShopkeeperDTO) error {
	v := validator.New()
	err := v.Struct(shopkeeper)

	if err != nil {
		return err
	}

	err = s.emailExists(shopkeeper.Email)
	if err != nil {
		return err
	}
	err = s.cnpjExists(shopkeeper.CNPJ)
	if err != nil {
		return err
	}

	return nil
}

func (s *Shopkeeper) cnpjExists(cpf string) error {
	_, err := s.client.Shopkeeper.FindUnique(db.Shopkeeper.Cnpj.Equals(cpf)).Exec(s.ctx)
	if err == nil {
		return errors.New("The field 'CNPJ' is already in use. Please use another CNPJ.")
	}
	return nil
}

func (s *Shopkeeper) getByEmail(email string) (*db.ShopkeeperModel, error) {
	shopkeeper, err := s.client.Shopkeeper.FindUnique(
		db.Shopkeeper.Email.Equals(email),
	).With(
		db.Shopkeeper.Wallet.Fetch(),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	return shopkeeper, nil
}

func MakeShopkeeperService(client *db.PrismaClient) *Shopkeeper {

	return &Shopkeeper{
		client: client,
	}
}
