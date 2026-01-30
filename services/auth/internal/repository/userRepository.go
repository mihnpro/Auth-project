package repository

import (
	"context"
	"database/sql"

	"github.com/mihnpro/Auth-project/services/auth/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, data *domain.User) (uint32, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

func (u *userRepository) CreateUser(ctx context.Context, data *domain.User) (uint32, error) {
	const query = `
	INSERT INTO users (email, password, phone_number)
	VALUES (?, ?, ?);`
	res, err := u.db.ExecContext(ctx, query, data.Email, data.Password, data.PhoneNumber)

	if err != nil {
		return 0, err
	}

	user_id, err := res.LastInsertId()

	return uint32(user_id), nil
}

func (u *userRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
	SELECT user_id, email, password, phone_number
	FROM users
	WHERE email = ?;`

	var user domain.User
	if err := u.db.QueryRowContext(ctx, query, email).Scan(&user.UserId, &user.Email, &user.Password, &user.PhoneNumber); err != nil {
		return nil, err
	}

	return &user, nil
}
