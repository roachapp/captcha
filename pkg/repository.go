package pkg

import (
	"context"
	"os"

	pgxc "github.com/jackc/pgx/v4"
	pgx "github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// Connect to postgres database
func ConnectDB(ctx context.Context) *pgx.Pool {
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

// PgNotifyer get the underlying connection of a pgxpool.Conn and opens a listening channel
// for notifications about events, see: https://github.com/jackc/pgx/issues/968
func PgNotifyer(ctx context.Context, pool *pgx.Pool) (*pgxc.Conn, error) {
	pgxpoolConn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("could not acquire postgres connection: %+v", err)
		return nil, err
	}

	notifyConn := pgxpoolConn.Conn()

	_, err = notifyConn.Exec(ctx, "LISTEN captchaNotifyChannel;") // TODO https://stackoverflow.com/questions/26046816/is-there-a-way-to-set-an-expiry-time-after-which-a-data-entry-is-automaticall
	if err != nil {
		log.Fatalf("could not establish SQL listener: %+v", err)
		return nil, err
	}

	return notifyConn, nil
}

// SelectCaptcha returns a PG transaction that queries a Captcha Rows
func SelectCaptcha() string {
	return "SELECT (solution, ttl) FROM captchas WHERE id = $1;"
}

// InsertCaptcha returns a PG transaction that creates an Captcha Row
func InsertCaptcha() string {
	return "INSERT INTO captchas (id, solution, ttl) VALUES ($1, $2, $3);"
}

// DeleteCaptcha returns a PG transaction that delete an Captcha Row by ID
func DeleteCaptcha() string {
	return "DELETE FROM captchas WHERE id = $1" // TODO https://stackoverflow.com/questions/26046816/is-there-a-way-to-set-an-expiry-time-after-which-a-data-entry-is-automaticall
}

