package validation

import (
	"github.com/klassmann/cpfcnpj"
	"gopkg.in/go-playground/validator.v9"
)

func validateCPF(fl validator.FieldLevel) bool {
	cpf := cpfcnpj.NewCPF(fl.Field().String())
	if !cpf.IsValid() {
		return false
	}

	return true
}
