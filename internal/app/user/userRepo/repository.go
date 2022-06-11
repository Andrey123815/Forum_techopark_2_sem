package userRepo

import (
	"db-forum/internal/app/models"
	"github.com/jackc/pgx"
)

type Repository struct {
	Database *pgx.ConnPool
}

func CreateUserRepository(db *pgx.ConnPool) *Repository {
	return &Repository{Database: db}
}

func (userRepo *Repository) GetUserByNickname(nickname string) (models.User, error) {
	var user models.User

	err := userRepo.Database.QueryRow(`SELECT nickname, fullname, about, email FROM users WHERE nickname = $1;`, nickname).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (userRepo *Repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := userRepo.Database.QueryRow(`SELECT nickname, fullname, about, email FROM users WHERE email = $1;`, email).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil && err.Error() == ("no rows in result set") {
		err = nil
	}
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (userRepo *Repository) CheckUserExists(nickname, email string) ([]models.User, error) {
	var users []models.User

	userWithSameNickname, err := userRepo.GetUserByNickname(nickname)
	if err != nil {
		userWithSameNickname = models.User{}
	} else {
		users = append(users, userWithSameNickname)
	}

	userWithSameEmail, err := userRepo.GetUserByEmail(email)
	if err != nil {
		userWithSameEmail = models.User{}
	}
	if userWithSameEmail.Email != "" && userWithSameNickname != userWithSameEmail {
		users = append(users, userWithSameEmail)
	}

	return users, nil
}

func (userRepo *Repository) CreateNewUser(nickname, fullname, about, email string) (models.User, error) {
	var newUser models.User

	err := userRepo.Database.QueryRow(`INSERT INTO users(nickname, fullname, about, email) VALUES ($1, $2, $3, $4)
		RETURNING nickname, fullname, about, email;`, nickname, fullname, about, email).
		Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (userRepo *Repository) UpdateUserProfile(nickname, fullname, about, email string) (models.User, error) {
	var newUser models.User

	err := userRepo.Database.QueryRow(`UPDATE users SET fullname=COALESCE(NULLIF($2, ''), fullname),
															about=COALESCE(NULLIF($3, ''), about),
															email=COALESCE(NULLIF($4, ''), email) WHERE nickname = $1
		RETURNING nickname, fullname, about, email;`, nickname, fullname, about, email).
		Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}
