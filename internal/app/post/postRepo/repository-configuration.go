package postRepo

import "db-forum/internal/app/models"

type PostRepository interface {
	GetPostDetails(postID int64, relatedArr []string) (map[string]interface{}, error)
	ChangePostMessage(postID int64, message string) (models.Post, error)
}
