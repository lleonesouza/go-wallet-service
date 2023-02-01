package dtos

type CreateShopkeeperDTO struct {
	Name     string `json:"name" validate:"required" example:"john"`
	Lastname string `json:"lastname" validate:"required" example:"doe"`
	CNPJ     string `json:"cnpj" validate:"required" example:"12345678"`
	Email    string `json:"email" validate:"required,email" example:"john@doe.com"`
	Password string `json:"password" validate:"required" example:"12345678"`
}

type UpdateShopkeeperDTO struct {
	Name     string `json:"name" validate:"required" example:"john"`
	Lastname string `json:"lastname" validate:"required" example:"doe"`
}

type LoginShopkeeperDTO struct {
	Email    string `json:"email" validate:"required,email" example:"john@doe.com"`
	Password string `json:"password" validate:"required,email" example:"12345678"`
}

//	@Description	Shopkeeper account information
type ShopkeeperResponseDTO struct {
	ID       string `json:"id" example:"06901d3b-134b-4ea6-ba0f-3a00ca5836b7"`
	Name     string `json:"name" example:"john"`
	Balance  int    `json:"balance" example:"50"`
	Lastname string `json:"lastname" example:"doe"`
	CNPJ     string `json:"cnpj" example:"12345789"`
	Email    string `json:"email" example:"jhon@doe.com"`
	CreateAt string `json:"create_at" example:"2023-01-31 12:47:27.072 +0000 UTC"`
	UpdateAt string `json:"update_at" example:"2023-01-31 12:47:27.072 +0000 UTC"`
}
