package userRepo

import (
	"db-forum/internal/app/models"
)

type UserRepository interface {
	GetUserByNickname(nickname string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CheckUserExists(nickname, email string) ([]models.User, error)

	CreateNewUser(nickname, fullname, about, email string) (models.User, error)
	UpdateUserProfile(nickname, fullname, about, email string) (models.User, error)
}
