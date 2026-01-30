package domain

type User struct {
	UserId      uint32 `db:"user_id"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	PhoneNumber string `db:"phone_number"`
}
