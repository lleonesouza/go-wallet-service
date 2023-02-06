package validation

import "gopkg.in/go-playground/validator.v9"

func MakeValidator() *validator.Validate {

	validate := validator.New()
	validate.RegisterValidation("cpf", validateCPF)
	validate.RegisterValidation("cnpj", validateCPF)
	validate.RegisterValidation("password", validatePassword)

	return validate
}
