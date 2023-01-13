package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net"
	"os"

	"google.golang.org/grpc"

	"ecomdream/src/contracts"
)


func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("IMAGER_PORT")))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxSendMsgSize(1024*1024*20),
		grpc.MaxRecvMsgSize(1024*1024*20),
	)

	contracts.RegisterImageServiceServer(s, &imageServiceServer{})

	if err := s.Serve(listener); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
