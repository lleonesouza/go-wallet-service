package dtos

type CreateShopkeeperDTO struct {
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
	CNPJ     string `json:"cnpj" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginShopkeeperDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,email"`
}

type UpdateShopkeeperDTO struct {
	Name     string `json:"name" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
}

type ShopkeeperResponseDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
	Lastname string `json:"lastname"`
	CNPJ     string `json:"cnpj"`
	Email    string `json:"email"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}
