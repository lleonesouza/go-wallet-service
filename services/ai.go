package services

import (
	"bff-answerfy/config"
	"bff-answerfy/prisma/db"
	"context"
	"net/http"
)

type Answer struct {
	Title  string
	Input  string
	Output string
}

type AI struct {
	client *db.PrismaClient
	http   *Http
	ctx    context.Context
	env    *config.Envs
}

func (ai *AI) SendChatGPT(text []byte) (*http.Response, error) {
	res, err := ai.http.POST(ai.env.OPENAI_URL, text)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ai *AI) GenerateAnswer(text string) (*ChatResponse, error) {
	res, err := ai.http.OpenAIRequest(text)
	if err != nil {
		return nil, err
	}

	return res, nil
}
