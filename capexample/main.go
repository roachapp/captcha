
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"time"

	"github.com/roachapp/captcha"
	pb "github.com/roachapp/captcha/api"
)

const (
	width = 160
	height = 80
)

type captchaServer struct {
	pb.UnimplementedCaptchaServer
	context context.Context
}

func (srv captchaServer) Validate(ctx context.Context, sol *pb.Solution) (*pb.Status, error) {
	if !captcha.VerifyString(sol.Id, sol.Code) {
		return &pb.Status{
			Code: 400,
			Message: "try again :(",
		}, nil
	}
	return &pb.Status{
		Code: 200,
		Message: "that went smoothly :)",
	}, nil

	// TODO send OK to backend servers containing users pubkey (pubkey = sol.Id)
}

func (srv captchaServer) Get(ctx context.Context, sol *pb.User) (*pb.Challenge, error) {
	captchaID := captcha.New()

	var content bytes.Buffer

	if err := captcha.WriteImage(&content, captchaID, width, height); err != nil {
		log.Error(err)
		return nil, err
	}

	return &pb.Challenge{
		Id: captchaID,
		Width: width,
		Height: height,
		GrayPixels: content.Bytes(),
	}, nil
}

type rateLimiter struct {
	rl *rate.Limiter
}

func (rl *rateLimiter) Limit() bool {
	return rl.rl.Allow()
}

func newServer() *grpc.Server {
	limiter := &rateLimiter{
		rl: rate.NewLimiter(rate.Every(time.Second * 30), 3), // note that a captcha's TTL is also 30 seconds
	}

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			ratelimit.UnaryServerInterceptor(limiter),
		),
	)
	pb.RegisterCaptchaServer(srv, captchaServer{context: context.Background()})
	return srv
}

func main() {
	ip := "0.0.0.0"
	port := 8666
	ipPort := fmt.Sprintf(ip + ":%d", port)

	// grpc approach
	conn, err := net.Listen("tcp", ipPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := newServer()
	log.Infof("Captcha Server running on %s", ipPort)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
