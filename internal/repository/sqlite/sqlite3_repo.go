package sqlite

import (
	"gorm.io/gorm"
)

type SqliteRepo struct {
	db *gorm.DB
}

type UserRepo struct {
	*SqliteRepo
}
