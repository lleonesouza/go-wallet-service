package dtos

type CreateTransactionDTO struct {
	Value int    `json:"value" validate:"required"`
	To    string `json:"to" validate:"required"`
}

type ResponseTransactionDTO struct {
	Id           string `json:"id"`
	FromWalletId string `json:"from_user_id"`
	ToWalletId   string `json:"to_user_id"`
	Value        int    `json:"value" `
	CreateAt     string `json:"create_at"`
	UpdateAt     string `json:"update_at"`
}
