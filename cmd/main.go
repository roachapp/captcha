
package main

import (
	"context"
	"fmt"
	"github.com/roachapp/captcha/pkg"
	log "github.com/sirupsen/logrus"
	"net"
	"os"

	"github.com/jackc/pgconn"
)

func main() {
	ip := "0.0.0.0"
	port := 8666
	ipPort := fmt.Sprintf(ip + ":%d", port)
	ctx := context.Background()

	// setup postgres database connection pool
	pg := pkg.ConnectDB(ctx)
	defer pg.Close()

	// setup postgres event notifyer connection
	pgNotifyer, err := pkg.PgNotifyer(ctx, pg)
	if err != nil {
		log.Fatalf("an error occured while establishing the event notifyer connection: %+v", err)
		os.Exit(1)
	}

	defer pgNotifyer.Close(ctx)

	notifyChan := make(chan *pgconn.Notification)

	go func() {
		// <-- upon SQL trigger, we notify the user of the events -->
		msg, err := pgNotifyer.WaitForNotification(ctx) // this blocks
		if err != nil {
			log.Errorf("error occured while waiting for SQL trigger: %+v", err)
		}

		// notify user
		notifyChan <- msg
		log.Debug("SQL trigger: triggered eventsNotifyChannel")
	}()

	// grpc approach
	conn, err := net.Listen("tcp", ipPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := pkg.NewServer(ctx)
	log.Infof("Captcha Server running on %s", ipPort)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
