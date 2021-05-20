package store

import (
	"context"
	"os"

	pgx "github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// Connect to postgres database
func connectDB(ctx context.Context) *pgx.Pool {
	dbURL, success := os.LookupEnv("DATABASE_URL")
	if !success {
		log.Fatalf("no environment-variable set for postgres DATABASE_URL")
		os.Exit(1)
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("could not connect to postgres database: %s", err)
		os.Exit(1)
	}

	log.Info("successfully connected to postgres database")
	return conn
}

// SelectCaptcha returns a PG transaction string that queries a Captcha Rows
func SelectCaptcha() string {
	return "SELECT (solution, pub_key) FROM captchas WHERE id = $1;"
}

// InsertCaptcha returns a PG transaction string that creates an Captcha Row
func InsertCaptcha() string {
	return "INSERT INTO captchas (id, solution) VALUES ($1, $2);"
}

// DeleteCaptcha returns a PG transaction string that deletes an Captcha Row by ID
func DeleteCaptcha() string {
	return "DELETE FROM captchas WHERE id = $1" // TODO https://stackoverflow.com/questions/26046816/is-there-a-way-to-set-an-expiry-time-after-which-a-data-entry-is-automaticall
}

