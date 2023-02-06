package services

import (
	"context"
	"q2bank/config"
	"q2bank/handlers/dtos"
	"q2bank/prisma/db"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Shopkeeper struct {
	client *db.PrismaClient
	wallet *Wallet
	ctx    context.Context
	env    *config.Envs
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

func (s *Shopkeeper) Create(_shopkeeper *dtos.CreateShopkeeperDTO, wallet *db.WalletModel) (*db.ShopkeeperModel, error) {
	passwordHash, err := HashPassword(_shopkeeper.Password)
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
	).With(db.Shopkeeper.Wallet.Fetch()).Update(
		db.Shopkeeper.Name.Set(_shopkeeper.Name),
		db.Shopkeeper.Lastname.Set(_shopkeeper.Lastname),
	).Exec(s.ctx)

	if err != nil {
		return shopkeeper, err
	}

	return shopkeeper, nil
}

func (s *Shopkeeper) Login(shopkeeper *db.ShopkeeperModel) (*string, error) {
	claims := &config.JwtCustomClaims{
		Email: shopkeeper.Email,
		ID:    shopkeeper.ID,
		Type:  s.env.SHOPKEEPER_TYPE,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.env.JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}

func (s *Shopkeeper) GetById(id string) (*db.ShopkeeperModel, error) {
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

func (s *Shopkeeper) GetByCPF(cpf string) (*db.UserModel, error) {
	shopkeeper, err := s.client.User.FindUnique(db.User.Cpf.Equals(cpf)).Exec(s.ctx)
	if err != nil {
		return nil, err
	}
	return shopkeeper, nil
}

func (s *Shopkeeper) GetByEmailWithWallet(email string) (*db.UserModel, error) {
	shopkeeper, err := s.client.User.FindUnique(db.User.Email.Equals(email)).With(db.User.Wallet.Fetch()).Exec(s.ctx)
	if err != nil {
		return nil, err
	}
	return shopkeeper, nil
}

func (s *Shopkeeper) GetByEmail(email string) (*db.UserModel, error) {
	shopkeeper, err := s.client.User.FindUnique(db.User.Email.Equals(email)).Exec(s.ctx)
	if err != nil {
		return nil, err
	}
	return shopkeeper, nil
}

func (s *Shopkeeper) CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func MakeShopkeeperService(client *db.PrismaClient) *Shopkeeper {

	return &Shopkeeper{
		client: client,
	}
}
