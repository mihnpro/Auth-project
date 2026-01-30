package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mihnpro/Auth-project/services/auth/internal/domain"
	"github.com/mihnpro/Auth-project/services/auth/internal/repository"
	"github.com/mihnpro/Auth-project/services/auth/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, data *UserCreateReq) (uint32, error)
	LoginUser(ctx context.Context, data *UserLoginReq) (string, string, uint32, error)
	RefreshTokens(ctx context.Context, refreshtoken string) (string, string, error)
}

type UserCreateReq struct {
	EmailAddress string
	Password     string
	PhoneNumber  string
}

type UserLoginReq struct {
	EmailAddress string
	Password     string
}


type userService struct {
	userRepo repository.UserRepository
	jwt      auth.JwtAuth
}

func NewUserService(userRepo repository.UserRepository, auth auth.JwtAuth) UserService {
	return &userService{
		userRepo: userRepo,
		jwt:      auth,
	}
}

func (u *userService) CreateUser(ctx context.Context, data *UserCreateReq) (uint32, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, data.EmailAddress)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if user != nil {
		return 0, errors.New("User has already exist")
	}

	hashedPassword, err := u.hashPassword(data.Password)

	if err != nil {
		return 0, errors.New("Error accured while hashing password")
	}

	userId, err := u.userRepo.CreateUser(ctx, &domain.User{
		Email:       data.EmailAddress,
		Password:    hashedPassword,
		PhoneNumber: data.PhoneNumber,
	})

	if userId == 0 {
		return 0, errors.New("Error accured while creating user")
	}

	return userId, nil

}

func (u *userService) hashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err

}

func (u *userService) LoginUser(ctx context.Context, data *UserLoginReq) (string, string, uint32, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, data.EmailAddress)

	if err != nil {
		return "", "", 0, err
	}

	if user == nil {
		return "", "", 0, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

	if err != nil {
		return "", "", 0, errors.New("Invalid password")
	}

	accessToken, refreshToken, err := u.jwt.GenerateTokens(user)

	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, user.UserId, nil

}

func (u *userService) RefreshTokens(ctx context.Context, refreshtoken string) (string, string, error) {
	newAccessToken, newRefreshToken, err := u.jwt.RefreshTokens(refreshtoken)

	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}
