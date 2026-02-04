package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mihnpro/Auth-project-protos/gen/go/permissionspb"
	"github.com/mihnpro/Auth-project/services/permissions/internal/repository"
	"github.com/mihnpro/Auth-project/services/permissions/internal/service"
	"github.com/mihnpro/Auth-project/services/permissions/internal/transport"
	"github.com/mihnpro/Auth-project/services/permissions/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	addr = "8081"
)

func main() {

	db := db.GetDB()

	permissionsRepo := repository.NewPermissionsRepository(db)

	permissionsService := service.NewPermissionsService(permissionsRepo)

	grpcServer := transport.NewGrpcServer(permissionsService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()

	permissionspb.RegisterPermissionsServiceServer(s, grpcServer)

	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-make(chan struct{})
}
