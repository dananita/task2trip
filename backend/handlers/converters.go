package handlers

import (
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/models"
)

func convertUser(user *backend.User) *models.User {
	return &models.User{
		ID:   &user.ID,
		Name: &user.Email,
	}
}
