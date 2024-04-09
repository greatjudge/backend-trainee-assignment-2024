// go:build integration

package tests

import (
	"banner/tests/postgres"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	db *postgres.TDB
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}

	psgDsn, ok := os.LookupEnv("POSTGRES_DB_DSN")
	if !ok {
		panic("no POSTGRES_DB_DSN in env vars")
	}

	// тут мы запрашиваем тестовые креды для бд из енв
	// cfg,err := config.FromEnv
	db = postgres.NewDB(psgDsn)
}
