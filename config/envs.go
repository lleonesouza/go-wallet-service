package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Envs struct {
	JWT_SECRET                 string
	TRANSACTION_AUTHORIZER_URL string
	APP_ENV                    string
	WALLET_BALANCE_INIT        string
	USER_TYPE                  string
	SHOPKEEPER_TYPE            string
}

func MakeEnvs() *Envs {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured loading env file. Err: %s", err)
	}

	envs := &Envs{
		JWT_SECRET:                 os.Getenv("JWT_SECRET"),
		TRANSACTION_AUTHORIZER_URL: os.Getenv("TRANSACTION_VERIFY_URL"),
		APP_ENV:                    os.Getenv("APP_ENV"),
		WALLET_BALANCE_INIT:        os.Getenv("WALLET_BALANCE_INIT"),
		USER_TYPE:                  "user",
		SHOPKEEPER_TYPE:            "shopkeeper",
	}

	return envs
}
