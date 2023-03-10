package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Envs struct {
	JWT_SECRET          string
	OPENAI_URL          string
	OPENAI_API_KEY      string
	APP_ENV             string
	WALLET_BALANCE_INIT string
	USER_TYPE           string
	SHOPKEEPER_TYPE     string
}

func MakeEnvs() *Envs {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured loading env file. Err: %s", err)
	}

	envs := &Envs{
		JWT_SECRET:          os.Getenv("JWT_SECRET"),
		OPENAI_URL:          "https://api.openai.com/v1/engine/chat",
		OPENAI_API_KEY:      os.Getenv("OPENAI_API_KEY"),
		APP_ENV:             os.Getenv("APP_ENV"),
		WALLET_BALANCE_INIT: os.Getenv("WALLET_BALANCE_INIT"),
		USER_TYPE:           "user",
	}

	return envs
}
