package business

import "hello/internal/domain"

type UserRepo interface {
	FindUsers(map[string]interface{}) (*domain.User, error)
	CreateUser(*domain.User)
}
