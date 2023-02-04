package validation

import (
	"github.com/klassmann/cpfcnpj"
	"gopkg.in/go-playground/validator.v9"
)

func validateCNPJ(fl validator.FieldLevel) bool {
	cnpj := cpfcnpj.NewCNPJ(fl.Field().String())
	if !cnpj.IsValid() {
		return false
	}

	return true

}
