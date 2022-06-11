package threadRepo

import (
	"db-forum/internal/app/models"
)

type ThreadRepository interface {
	GetThreadByID(threadID int64) (models.Thread, error)
	GetThreadIDBySlugOrID(slugOrID string) (int64, error)
	GetThreadBySlugOrID(slugOrID string) (models.Thread, error)

	CreateNewPosts(newPosts []models.Post, thread models.Thread) ([]models.Post, error)
	VoteThread(nickname string, threadID int64, voice int32, oldTread models.Thread) (models.Thread, error)
	GetThreadPosts(id int64, limit, since, sortType, sortDirection string) ([]models.Post, error)
	UpdateThreadDetails(threadID int64, title, message string) (models.Thread, error)
}
