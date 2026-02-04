package transport

import (
	permissionspb "github.com/mihnpro/Auth-project-protos/gen/go/permissionspb"
	"github.com/mihnpro/Auth-project/services/permissions/internal/service"
)

type GrpcServer struct {
	permissionspb.UnimplementedPermissionsServiceServer
	service service.PermissionsService
}


func NewGrpcServer(service service.PermissionsService) *GrpcServer {
	return &GrpcServer{
		service: service,
	}
}

