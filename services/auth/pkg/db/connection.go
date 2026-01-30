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
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		phone_number TEXT NOT NULL UNIQUE
	);`
	
	if _, err := DB.ExecContext(ctx, query); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Table 'users' initialized/checked")
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