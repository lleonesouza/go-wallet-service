package services

import (
	"bff-answerfy/config"
	"bff-answerfy/handlers/dtos"
	"bff-answerfy/prisma/db"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	client *db.PrismaClient
	wallet *Wallet
	ctx    context.Context
	env    *config.Envs
}

func (u *User) Create(_user *dtos.CreateUserDTO, wallet *db.WalletModel) (*db.UserModel, error) {
	passwordHash, err := HashPassword(_user.Password)
	if err != nil {
		return nil, err
	}

	user, err := u.client.User.CreateOne(
		db.User.Email.Set(_user.Email),
		db.User.Cpf.Set(_user.CPF),
		db.User.Password.Set(string(passwordHash)),
		db.User.Name.Set(_user.Name),
		db.User.Lastname.Set(_user.Lastname),
		db.User.Wallet.Link(
			db.Wallet.ID.Equals(wallet.ID),
		),
	).With(db.User.Wallet.Fetch()).Exec(u.ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Update(id string, _user *dtos.UpdateUserDTO) (*db.UserModel, error) {
	user, err := u.client.User.FindUnique(
		db.User.ID.Equals(id),
	).With(db.User.Wallet.Fetch()).Update(
		db.User.Name.Set(_user.Name),
		db.User.Lastname.Set(_user.Lastname),
	).Exec(u.ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Login(user *db.UserModel) (*string, error) {
	claims := &config.JwtCustomClaims{
		Email: user.Email,
		ID:    user.ID,
		Type:  u.env.USER_TYPE,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(u.env.JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}

func (u *User) GetById(id string) (*db.UserModel, error) {
	user, err := u.client.User.FindUnique(
		db.User.ID.Equals(id),
	).With(
		db.User.Wallet.Fetch(),
	).Exec(u.ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetByCPF(cpf string) (*db.UserModel, error) {
	user, err := u.client.User.FindUnique(db.User.Cpf.Equals(cpf)).Exec(u.ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetByEmailWithWallet(email string) (*db.UserModel, error) {
	user, err := u.client.User.FindUnique(
		db.User.Email.Equals(email),
	).With(
		db.User.Wallet.Fetch(),
	).Exec(u.ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetByEmail(email string) (*db.UserModel, error) {
	user, err := u.client.User.FindUnique(db.User.Email.Equals(email)).Exec(u.ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Filter(user *db.UserModel) *dtos.UserResponseDTO {
	return &dtos.UserResponseDTO{
		ID:       user.ID,
		WalletID: user.WalletID,
		Balance:  user.Wallet().Balance,
		Name:     user.Name,
		Lastname: user.Lastname,
		CPF:      user.Cpf,
		Email:    user.Email,
		CreateAt: user.CreatedAt.String(),
		UpdateAt: user.UpdatedAt.String(),
	}
}

func (u *User) CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
