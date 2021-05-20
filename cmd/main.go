
package main

import (
	"context"
	"fmt"
	"github.com/roachapp/captcha/pkg/captcha"
	"github.com/roachapp/captcha/pkg/store"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func main() {
	ip := "0.0.0.0"
	port := 8666
	ipPort := fmt.Sprintf(ip + ":%d", port)
	ctx := context.Background()

	// create captcha generator
	captchaGenerator := &captcha.Generator{
		DigitLen: 3,
		Width: 160,
		Height: 80,
		CacheStore: store.NewCacheStore(100, 30 * time.Second),
		PgStore:    store.NewPostgresStore(ctx),
	}

	// grpc connection
	conn, err := net.Listen("tcp", ipPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := captcha.NewServer(ctx, captchaGenerator)
	log.Infof("Captcha Server running on %s", ipPort)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
