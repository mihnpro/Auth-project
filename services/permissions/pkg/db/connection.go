package db

import (
	"context"
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB   *sql.DB
	once sync.Once
)

func InitDB() *sql.DB {
	once.Do(func() {
		connection()
		initTable(context.Background())
	})
	return DB
}

func GetDB() *sql.DB {
	if DB == nil {
		return InitDB()
	}
	return DB
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

func initTable(ctx context.Context) {
	const query = `
	CREATE TABLE IF NOT EXISTS permissions (
		permission_id INTEGER PRIMARY KEY AUTOINCREMENT,
		permission_name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS users_with_permissions (
		user_id INTEGER NOT NULL,
		permission_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (permission_id) REFERENCES permissions(permission_id)
	);
	`

	if _, err := DB.ExecContext(ctx, query); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Tables 'permissions' and 'users_with_permissions' initialized/checked")
}

func connection() {
	var err error
	DB, err = sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Проверяем соединение
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Настройки пула соединений
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	log.Println("Database connection established")
}
