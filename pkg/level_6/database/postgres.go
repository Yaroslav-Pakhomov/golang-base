package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

// CheckConnect - проверка соединения с БД
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

// CreatePostsTable - создание табл. Поста
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

// CreatePost - создание поста
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

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// SelectPosts - получение всех постов
func SelectPosts(db *sql.DB) ([]Post, error) {

	rows, err := db.Query("SELECT id, title, description, sort_order, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.SortOrder, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostById - получение Поста по ID
func GetPostById(db *sql.DB, id int) (Post, error) {
	var post Post

	row := db.QueryRow(
		"SELECT id, title, description, sort_order, created_at FROM posts WHERE id = $1",
		id,
	)

	err := row.Scan(&post.ID, &post.Title, &post.Description, &post.SortOrder, &post.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Post{}, fmt.Errorf("post not found")
	}

	if err != nil {
		return Post{}, err
	}

	return post, nil
}
