package transport

import (
	authpb "github.com/mihnpro/Auth-project-protos/gen/go/authpb"
	"github.com/mihnpro/Auth-project/services/auth/internal/service"
)

type GrpcServer struct {
	authpb.UnimplementedAuthServiceServer
	service service.UserService
}


func NewGrpcServer(service service.UserService) *GrpcServer {
	return &GrpcServer{
		service: service,
	}
}

