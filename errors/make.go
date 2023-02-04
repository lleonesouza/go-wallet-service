package errors

import (
	"fmt"
	"net/http"
)

type Errors struct {
	BodyError         func(details string) *Error
	EmailRegistered   func(email string) *Error
	CpfRegistered     func(cpf string) *Error
	TypeError         func(details string) *Error
	UnauthorizedError func() *Error
	CreateWalletError func(details string) *Error
	CreateUserError   func(details string) *Error
	UpdateUserError   func(details string) *Error
}

type Error struct {
	Title   string `json:"title" example:"body_error"`
	Details string `json:"details" example:"Key: 'CreateUserDTO.Password' Error:Field validation for 'Password' failed on the 'password' tag\nKey: 'CreateUserDTO.CPF' Error:Field validation for 'CPF' failed on the 'cpf' tag"`
	Status  int    `json:"status" example:"422"`
}

func MakeErrors() *Errors {

	return &Errors{
		BodyError: func(details string) *Error {
			return &Error{
				Title:   "body_error",
				Details: details,
				Status:  http.StatusUnprocessableEntity,
			}
		},
		EmailRegistered: func(email string) *Error {
			return &Error{
				Title:   "email_conflict",
				Details: fmt.Sprintf("The email '%s' is already registered.", email),
				Status:  http.StatusConflict,
			}
		},
		CpfRegistered: func(cpf string) *Error {
			return &Error{
				Title:   "cpf_conflict",
				Details: fmt.Sprintf("The cpf '%s' is already registered.", cpf),
				Status:  http.StatusConflict,
			}
		},
		TypeError: func(details string) *Error {
			return &Error{
				Title:   "type_error",
				Details: "Body type is incorrect. Please check the fields or visit or docs: /swagger",
				Status:  http.StatusBadRequest,
			}
		},
		CreateWalletError: func(details string) *Error {
			return &Error{
				Title:   "create_wallet_error",
				Details: "Fails trying to create Wallet, please try again",
				Status:  http.StatusBadRequest,
			}
		},
		CreateUserError: func(details string) *Error {
			return &Error{
				Title:   "create_user_error",
				Details: "Fails trying to create User, please try again",
				Status:  http.StatusBadRequest,
			}
		},
		UnauthorizedError: func() *Error {
			return &Error{
				Title:   "unauthorized_error",
				Details: "Please log-in with your account",
				Status:  http.StatusUnauthorized,
			}
		},
	}
}
