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

type JwtCustomClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}

type User struct {
	client *db.PrismaClient
	wallet *Wallet
	ctx    context.Context
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

func (u *User) Create(_user *dtos.CreateUserDTO) (*db.UserModel, error) {
	err := u.validate(_user)
	if err != nil {
		return nil, err
	}

	passwordHash, err := HashPassword(_user.Password)
	if err != nil {
		return nil, err
	}

	wallet, err := u.wallet.Create()
	if err != nil {
		return nil, err
	}

	user, err := u.client.User.CreateOne(
		db.User.Email.Set(_user.Email),
		db.User.Cpf.Set(_user.CPF),
		db.User.Password.Set(passwordHash),
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
	).Update(
		db.User.Name.Set(_user.Name),
		db.User.Lastname.Set(_user.Lastname),
	).Exec(u.ctx)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) Get(id string) (*db.UserModel, error) {
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

func (u *User) Login(email string, password string) (string, error) {
	user, err := u.getByEmail(email)
	if err != nil {
		return "", errors.New("Email or Password incorrect")
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("Email or Password incorrect")
	}

	claims := &config.JwtCustomClaims{
		Email: user.Email,
		ID:    user.ID,
		Type:  "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))

	return t, nil
}

func (u *User) emailExists(email string) error {
	_, err := u.client.User.FindUnique(db.User.Email.Equals(email)).Exec(u.ctx)
	if err == nil {
		return errors.New("The field 'Email' is already in use. Please use another email.")
	}
	return nil
}

func (u *User) cpfExists(cpf string) error {
	_, err := u.client.User.FindUnique(db.User.Cpf.Equals(cpf)).Exec(u.ctx)
	if err == nil {
		return errors.New("The field 'CPF' is already in use. Please use another CPF.")
	}
	return nil
}

func (u *User) getByEmail(email string) (*db.UserModel, error) {
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

func (u *User) validate(user *dtos.CreateUserDTO) error {
	v := validator.New()
	err := v.Struct(user)

	if err != nil {
		return err
	}

	err = u.emailExists(user.Email)
	if err != nil {
		return err
	}
	err = u.cpfExists(user.CPF)
	if err != nil {
		return err
	}

	return nil
}

func MakeUserService(client *db.PrismaClient) *User {
	return &User{
		client: client,
	}
}
