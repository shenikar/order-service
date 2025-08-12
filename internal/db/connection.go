package db

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/shenikar/order-service/config"
)

// Connect устанавливает соединение с базой данных
func Connect(config *config.Config) (*sqlx.DB, error) {
	connStr := config.GetDatabaseUrl()

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to the database successfully")
	return db, nil
}
