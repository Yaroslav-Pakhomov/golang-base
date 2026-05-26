package database

import (
	"context"
	"database/sql"
	"fmt"

	// Конфиг приложения с данными для подключения к БД
	"golang-base/pkg/level_6/config"

	// Регистрируем драйвер pgx для database/sql.
	// Через "_" импортируется только init() драйвера.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ConnectPostgresDb - создаёт подключение к PostgreSQL
// и проверяет его через PingContext.
func ConnectPostgresDb(ctx context.Context, cfg *config.Config) (*sql.DB, error) {

	// Формируем DSN (Data Source Name) —
	// строку подключения к PostgreSQL.
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
		cfg.DBSSLMode,
	)

	// Создаём объект подключения к БД.
	//
	// Важно:
	// sql.Open НЕ открывает соединение сразу,
	// а только подготавливает пул соединений.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем реальное подключение к PostgreSQL.
	//
	// PingContext:
	// - пытается установить соединение
	// - учитывает context.Context
	// - позволяет использовать timeout/cancel
	if err = db.PingContext(ctx); err != nil {

		// Если подключение не удалось —
		// закрываем пул соединений.
		db.Close()
		return nil, err
	}

	// Возвращаем готовое подключение к БД.
	return db, nil
}

func CheckConnect(ctx context.Context, db *sql.DB) error {
	var version string

	err := db.QueryRowContext(
		ctx,
		"SELECT version()",
	).Scan(&version)

	if err != nil {
		return err
	}

	fmt.Println(version)

	return nil
}

func CreatePostsTable(db *sql.DB) error {
	query := `
			CREATE TABLE IF NOT EXISTS posts (
				id SERIAL PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT NOT NULL,
				sort_order INT NOT NULL,
				created_at TIMESTAMP DEFAULT NOW()
			);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func CreatePost(db *sql.DB, title string, description string, sort_order int) error {
	_, err := db.Exec(
		"INSERT INTO posts (title, description, sort_order) VALUES ($1, $2, $3)",
		title,
		description,
		sort_order,
	)

	if err != nil {
		return err
	}

	return nil
}
