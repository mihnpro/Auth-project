package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mihnpro/Auth-project/services/auth/pkg/auth"
	"github.com/mihnpro/Auth-project-protos/gen/go/authpb"
	"github.com/mihnpro/Auth-project/services/auth/internal/repository"
	"github.com/mihnpro/Auth-project/services/auth/internal/service"
	"github.com/mihnpro/Auth-project/services/auth/internal/transport"
	"github.com/mihnpro/Auth-project/services/auth/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	addr = "8080"
)

func main() {

	db := db.GetDB()

	userRepo := repository.NewUserRepository(db)

	jwtService := auth.NewJwtService()

	userService :=	service.NewUserService(userRepo, jwtService)

	grpcServer := transport.NewGrpcServer(userService)


	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()

	authpb.RegisterAuthServiceServer(s, grpcServer)

	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	select {}
}
