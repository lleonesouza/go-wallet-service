package dtos

type CreateQuestionDTO struct {
	Text string `json:"text" validate:"required" example:"O que é biologia celular?"`
}
