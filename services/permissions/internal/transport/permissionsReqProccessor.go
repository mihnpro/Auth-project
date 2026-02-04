package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mihnpro/Auth-project-protos/gen/go/permissionspb"
	"github.com/mihnpro/Auth-project/services/permissions/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GrpcServer) AssignRole(ctx context.Context, req *permissionspb.AssignRoleRequest) (*emptypb.Empty, error) {
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(time.RFC822))
		grpc.SetTrailer(ctx, trailer)
	}()

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		if vals, ok := md["timestamp"]; ok && len(vals) > 0 {
			fmt.Printf("timestamp metadata:")
			for i, v := range vals {
				fmt.Printf(" %d: %s;", i, v)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("no metadata received")
	}

	if !ok {
		log.Println("No metadata found")
	}

	err := g.service.AssignRole(ctx, &service.ProccessPermissionsReq{
		UserId:   req.UserId,
		RoleName: req.Role,
	})

	if err != nil {
		return nil, err
	}

	header := metadata.New(map[string]string{"location": "SCH", "timestamp": time.Now().Format(time.RFC822)})

	grpc.SetHeader(ctx, header)

	fmt.Println("request recived")

	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) CheckPermissions(ctx context.Context, req *permissionspb.CheckPermissionRequest) (*permissionspb.CheckPermissionResponse, error) {
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(time.RFC822))
		grpc.SetTrailer(ctx, trailer)
	}()

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		if vals, ok := md["timestamp"]; ok && len(vals) > 0 {
			fmt.Printf("timestamp metadata:")
			for i, v := range vals {
				fmt.Printf(" %d: %s;", i, v)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("no metadata received")
	}

	if !ok {
		log.Println("No metadata found")
	}

	allowed, err := g.service.CheckPermissions(ctx, &service.ProccessPermissionsReq{
		UserId:   req.UserId,
		RoleName: req.Permission,
	})

	if err != nil {
		return nil, err
	}

	header := metadata.New(map[string]string{"location": "SCH", "timestamp": time.Now().Format(time.RFC822)})

	grpc.SetHeader(ctx, header)

	fmt.Println("request recived")

	return &permissionspb.CheckPermissionResponse{
		Allowed: allowed,
	}, nil
}

func (g *GrpcServer) GetUserPermissions(ctx context.Context, req *permissionspb.GetUserPermissionsRequest) (*permissionspb.GetUserPermissionsResponse, error) {
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(time.RFC822))
		grpc.SetTrailer(ctx, trailer)
	}()

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		if vals, ok := md["timestamp"]; ok && len(vals) > 0 {
			fmt.Printf("timestamp metadata:")
			for i, v := range vals {
				fmt.Printf(" %d: %s;", i, v)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("no metadata received")
	}

	if !ok {
		log.Println("No metadata found")
	}

	permissions, err := g.service.GetPermissions(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	if permissions == "" {
		return nil, errors.New("No permissions")
	}

	header := metadata.New(map[string]string{"location": "SCH", "timestamp": time.Now().Format(time.RFC822)})

	grpc.SetHeader(ctx, header)

	fmt.Println("request recived")

	return &permissionspb.GetUserPermissionsResponse{
		Permissions: permissions,
	}, nil
}
