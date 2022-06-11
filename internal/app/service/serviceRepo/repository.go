package serviceRepo

import (
	"db-forum/internal/app/models"
	"github.com/jackc/pgx"
)

type Repository struct {
	Database *pgx.ConnPool
}

func CreateServiceRepository(db *pgx.ConnPool) *Repository {
	return &Repository{Database: db}
}

func (repository *Repository) GetSystemStatus() models.Service {
	var serviceInfo models.Service

	_ = repository.Database.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&serviceInfo.User)
	_ = repository.Database.QueryRow(`SELECT COUNT(*) FROM forums;`).Scan(&serviceInfo.Forum)
	_ = repository.Database.QueryRow(`SELECT COUNT(*) FROM threads;`).Scan(&serviceInfo.Thread)
	_ = repository.Database.QueryRow(`SELECT COUNT(*) FROM posts;`).Scan(&serviceInfo.Post)

	return serviceInfo
}

func (repository *Repository) ClearSystem() error {
	err := repository.Database.QueryRow(`TRUNCATE users, forums, posts, threads, votes, forum_users CASCADE`).Scan()
	return err
}
