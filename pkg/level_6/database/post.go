package database

import (
	"errors"
	"time"
)

// ErrPostNotFound — пост с указанным id не найден в БД.
// Используется и в postgres.go (database/sql), и в postgresGorm.go (GORM).
var ErrPostNotFound = errors.New("post not found")

// Post — модель поста для всего пакета database.
//
// Один тип — два набора тегов:
//   - json  — сериализация в API (ответы HTTP, логи)
//   - gorm  — маппинг на таблицу posts в PostgreSQL
//
// Как config.Config вынесен в config.go, так Post вынесен сюда,
// чтобы postgres.go и postgresGorm.go не дублировали структуру.
type Post struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	SortOrder   int       `json:"sort_order" gorm:"column:sort_order;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

// TableName — имя таблицы для GORM.
// Без этого метода GORM искал бы таблицу "posts" по умолчанию для Post — совпадает,
// но явное имя фиксирует контракт с CREATE TABLE в postgres.go.
func (Post) TableName() string {
	return "posts"
}
