package services

import (
	"bff-answerfy/config"
	"bff-answerfy/prisma/db"
	"context"
)

type Question struct {
	client *db.PrismaClient
	wallet *Wallet
	ctx    context.Context
	env    *config.Envs
}

func (q *Question) Create(user_id string, text string) (*db.QuestionModel, error) {
	question, err := q.client.Question.CreateOne(
		db.Question.Status.Set("processing"),
		db.Question.Input.Set(text),
		db.Question.User.Link(
			db.User.ID.Equals(user_id),
		),
	).Exec(q.ctx)

	if err != nil {
		return nil, err
	}

	return question, nil
}

func (q *Question) Update(question_id string, text string, status string) (*db.QuestionModel, error) {
	question, err := q.client.Question.FindUnique(
		db.Question.ID.Equals(question_id),
	).Update(
		db.Question.Output.Set(text),
		db.Question.Status.Set(status),
	).Exec(q.ctx)

	if err != nil {
		return nil, err
	}

	return question, nil
}

func (q *Question) GetAll(user_id string) (*[]db.QuestionModel, error) {
	questions, err := q.client.Question.FindMany(
		db.Question.UserID.Equals(user_id),
	).Exec(q.ctx)

	if err != nil {
		return nil, err
	}

	return &questions, nil
}
