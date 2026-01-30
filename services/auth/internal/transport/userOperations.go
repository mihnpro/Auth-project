package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	authpb "github.com/mihnpro/Auth-project-protos/gen/go/authpb"
	"github.com/mihnpro/Auth-project/services/auth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (g *GrpcServer) SignUp(ctx context.Context, req *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {

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

	user_id, err := g.service.CreateUser(ctx, &service.UserCreateReq{
		EmailAddress: req.EmailAddress,
		Password:     req.Password,
		PhoneNumber:  req.PhoneNumber,
	})

	if err != nil {
		return nil, err
	}

	if user_id == 0 {
		return nil, err
	}

	if !ok {
		log.Println("No metadata found")
	}

	header := metadata.New(map[string]string{"location": "SCH", "timestamp": time.Now().Format(time.RFC822)})

	grpc.SetHeader(ctx, header)

	fmt.Println("request recived")

	return &authpb.SignUpResponse{Message: "User created",
		UserId: uint32(user_id)}, nil

}

func (g *GrpcServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
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

	access_token, refresh_token, user_id, err := g.service.LoginUser(ctx, &service.UserLoginReq{
		EmailAddress: req.EmailAddress,
		Password:     req.Password,
	})

	if err != nil {
		return nil, err
	}

	if user_id == 0 {
		return nil, err
	}

	if access_token == "" {
		return nil, errors.New("No access token")
	}

	if refresh_token == "" {
		return nil, errors.New("No refresh token")
	}

	if !ok {
		log.Println("No metadata found")
	}

	header := metadata.New(map[string]string{"location": "SCH", "timestamp": time.Now().Format(time.RFC822)})

	grpc.SetHeader(ctx, header)

	fmt.Println("request recived")

	return &authpb.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		UserId:       uint32(user_id),
	}, nil
}

func (g *GrpcServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	return nil, nil
}

func (g *GrpcServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
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

	access_token, refresh_token, err := g.service.RefreshTokens(ctx, req.RefreshToken)

	if err != nil {
		return nil, err
	}

	if access_token == "" {
		return nil, errors.New("No access token")
	}

	if refresh_token == "" {
		return nil, errors.New("No refresh token")
	}

	if !ok {
		log.Println("No metadata found")
	}

	header := metadata.New(map[string]string{"location": "SCH", "timestamp": time.Now().Format(time.RFC822)})

	grpc.SetHeader(ctx, header)

	fmt.Println("request recived")

	return &authpb.RefreshTokenResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}
