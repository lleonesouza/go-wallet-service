package dtos

// Inputs

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
	CPF      string `json:"cpf" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,email"`
}

// Outputs
type UserResponseDTO struct {
	ID       string `json:"id"`
	WalletID string `json:"wallet_id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Balance  int    `json:"balance"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type LoginResponseDTO struct {
	Token int `json:"token"`
}
