package forumRepo

import (
	"db-forum/internal/app/models"
	"time"
)

type ForumRepository interface {
	GetForumBySlug(slug string) (models.Forum, error)
	GetAlreadyExistThread(slug string) (models.Thread, error)
	GetThreadsBySlug(forumSlug string, limit string, desc bool, since string) ([]models.Thread, error)

	CreateNewForum(title, user, slug string) (models.Forum, error)
	CreateNewThread(title, author, forum, message, slug string, created time.Time) (models.Thread, error)
	GetForumUsers(forumID int64, limit string, desc string, since string) ([]models.User, error)
}
