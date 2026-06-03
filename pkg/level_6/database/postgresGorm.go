package database

import (
	"context"
	"errors"
	"fmt"

	// Конфиг приложения с данными для подключения к БД.
	"golang-base/pkg/level_6/config"

	// Драйвер PostgreSQL для GORM (внутри использует pgx).
	"gorm.io/driver/postgres"

	// Ядро GORM — ORM (Object-Relational Mapping):
	// вместо ручного SQL мы работаем со структурами Go.
	"gorm.io/gorm"
)

// ConnectPostgresGormDb — создаёт подключение к PostgreSQL через GORM
// и проверяет его через PingContext.
//
// Аналог ConnectPostgresDb в postgres.go, но возвращает *gorm.DB вместо *sql.DB.
func ConnectPostgresGormDb(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {

	// DSN для GORM-драйвера postgres — формат key=value (не URL).
	// В postgres.go используется URL: postgres://user:pass@host:port/db
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBDatabase,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	// gorm.Open:
	// - postgres.Open(dsn) — настраивает драйвер
	// - &gorm.Config{} — настройки ORM (логирование, именование и т.д.)
	//
	// Как и sql.Open, соединение с БД может не установиться сразу —
	// поэтому ниже делаем Ping.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// db.DB() возвращает нижележащий *sql.DB из пакета database/sql.
	// GORM построен поверх стандартного пула соединений.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// PingContext — та же проверка «живости» БД, что в postgres.go.
	if err = sqlDB.PingContext(ctx); err != nil {

		// При ошибке закрываем пул; _ игнорирует ошибку Close.
		_ = sqlDB.Close()
		return nil, err
	}

	return db, nil
}

// CheckConnectGorm — проверка соединения с БД через GORM.
//
// Здесь намеренно используем «сырой» *sql.DB:
// GORM не обязателен для простого SELECT version().
func CheckConnectGorm(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	var version string

	err = sqlDB.QueryRowContext(
		ctx,
		"SELECT version()",
	).Scan(&version)

	if err != nil {
		return err
	}

	fmt.Println(version)

	return nil
}

// CreatePostsTableGorm — создание/обновление схемы таблицы posts.
//
// AutoMigrate сравнивает модель Post с таблицей в БД и при необходимости:
// - создаёт таблицу, если её нет
// - добавляет недостающие колонки
//
// В postgres.go то же делается вручную через CREATE TABLE IF NOT EXISTS.
func CreatePostsTableGorm(db *gorm.DB) error {
	return db.AutoMigrate(&Post{})
}

// CreatePostGorm — создание поста (INSERT).
//
// В postgres.go: Begin → Exec(INSERT) → Commit.
// GORM сам формирует INSERT и по умолчанию выполняет его в одной операции.
func CreatePostGorm(db *gorm.DB, title string, description string, sortOrder int) error {

	// Заполняем только поля, которые передаём из кода.
	// ID и CreatedAt GORM/БД заполнят автоматически.
	post := Post{
		Title:       title,
		Description: description,
		SortOrder:   sortOrder,
	}

	// Create(&post):
	// - генерирует SQL INSERT
	// - после успеха записывает сгенерированный id обратно в post.ID
	//
	// .Error — в GORM ошибки лежат в result.Error, не возвращаются вторым значением.
	return db.Create(&post).Error
}

// SelectPostsGorm — получение всех постов (SELECT).
//
// Аналог SelectPosts: вместо rows.Next() + Scan() — один вызов Find.
func SelectPostsGorm(db *gorm.DB) ([]Post, error) {
	var posts []Post

	// Find(&posts) ≈ SELECT * FROM posts ORDER BY id
	// Order("id") — сортировка как при чтении по порядку id.
	if err := db.Order("id").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostByIdGorm — получение одного поста по ID.
//
// First(&post, id) ≈ SELECT * FROM posts WHERE id = ? LIMIT 1
func GetPostByIdGorm(db *gorm.DB, id int) (Post, error) {
	var post Post

	err := db.First(&post, id).Error

	// gorm.ErrRecordNotFound — аналог sql.ErrNoRows в postgres.go.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Post{}, ErrPostNotFound
	}

	if err != nil {
		return Post{}, err
	}

	return post, nil
}

// UpdatePostGorm — обновление поста по ID.
//
// В postgres.go проверяют RowsAffected после UPDATE.
// Здесь то же: если строк с таким id нет — ErrPostNotFound.
func UpdatePostGorm(db *gorm.DB, id int, title string, description string, sortOrder int) error {

	// Model(&Post{}) — указываем таблицу posts.
	// Where("id = ?", id) — условие; ? подставляется безопасно (не конкатенация строк).
	// Updates(map) — обновляет только перечисленные колонки.
	result := db.Model(&Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":       title,
		"description": description,
		"sort_order":  sortOrder,
	})

	if result.Error != nil {
		return result.Error
	}

	// RowsAffected — сколько строк изменилось (0 = пост не найден).
	if result.RowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}

// DeletePostGorm — удаление поста по ID.
//
// Delete(&Post{}, id) ≈ DELETE FROM posts WHERE id = ?
func DeletePostGorm(db *gorm.DB, id int) error {
	result := db.Delete(&Post{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}
