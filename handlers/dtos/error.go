package dtos

type EmailExistsError struct {
	Error string `json:"error" example:"email already exists"`
}

type CPFExistsError struct {
	Error string `json:"error" example:"cpf already exists"`
}

type UnauthorizedError struct {
	Error string `json:"error" example:"missing or malformed jwt"`
}

type GeneralError struct {
	Error string `json:"error" example:"something goes wrong"`
}
