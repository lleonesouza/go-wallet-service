package dtos

// Inputs

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required" example:"john"`
	Lastname string `json:"lastname" validate:"required" example:"doe"`
	CPF      string `json:"cpf" validate:"required" example:"12345678"`
	Email    string `json:"email" validate:"required,email" example:"john@doe.com"`
	Password string `json:"password" validate:"required" example:"12345678"`
}

type UpdateUserDTO struct {
	Name     string `json:"name" validate:"required" example:"john"`
	Lastname string `json:"lastname" validate:"required" example:"doe"`
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email" example:"john@doe.com"`
	Password string `json:"password" validate:"required,email" example:"12345678"`
}

// Outputs
type UserResponseDTO struct {
	ID       string `json:"id" example:"06901d3b-134b-4ea6-ba0f-3a00ca5836b7"`
	WalletID string `json:"wallet_id" example:"06901d3b-134b-4ea6-ba0f-3a00ca5836b7"`
	Name     string `json:"name" example:"john"`
	Lastname string `json:"lastname" example:"doe"`
	CPF      string `json:"cpf" example:"12345678"`
	Email    string `json:"email" example:"john@doe.com"`
	Balance  int    `json:"balance" example:"50"`
	CreateAt string `json:"create_at" example:"2023-01-31 12:47:27.072 +0000 UTC"`
	UpdateAt string `json:"update_at" example:"2023-01-31 12:47:27.072 +0000 UTC"`
}

type LoginResponseDTO struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZG9lLmNvbSIsImlkIjoiMWE4MjQwM2YtYWNhOS00YjA1LTliNTEtYjRmZWE4OGM2MWQ5IiwidHlwZSI6InNob3BrZWVwZXIiLCJleHAiOjE2NzU1NDkyODd9.MSgwyCzvbC6tfH7ZYNrEhhv_XbmKqVEX-rEe6Y7EMKI"`
}
