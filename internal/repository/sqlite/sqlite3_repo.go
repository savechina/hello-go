package sqlite

import (
	"hello/internal/domain"

	"gorm.io/gorm"
)

type SqliteRepo struct {
	db *gorm.DB
}

type UserRepDao interface {
	FindUsers(map[string]interface{}) (*domain.User, error)
	CreateUser(*domain.User)
}

type UserRepo struct {
	*SqliteRepo
}
