package services

import "github/fabulousCode/models"

type UserService interface {
	CreateUser(*models.User) error
	GetAllUsers() ([]*models.User, error)
	GetUser(*string) (*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}
