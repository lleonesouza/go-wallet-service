package dtos

type CreateTransactionDTO struct {
	Value int    `json:"value" validate:"required" example:"20"`
	To    string `json:"to" validate:"required" example:"john@doe.com"`
}

type ResponseTransactionDTO struct {
	Id           string `json:"id" example:"06901d3b-134b-4ea6-ba0f-3a00ca5836b7"`
	FromWalletId string `json:"from_user_id" example:"06901d3b-134b-4ea6-ba0f-3a00ca5836b7"`
	ToWalletId   string `json:"to_user_id" example:"06901d3b-134b-4ea6-ba0f-3a00ca5836b7"`
	Value        int    `json:"value" example:"50"`
	CreateAt     string `json:"create_at" example:"2023-01-31 12:47:27.072 +0000 UTC"`
	UpdateAt     string `json:"update_at" example:"2023-01-31 12:47:27.072 +0000 UTC"`
}
