package postRepo

import "db-forum/internal/app/models"

type PostRepository interface {
	GetPostDetails(id int, related []string) (postInfo models.DetailPostInfo, err error)
	ChangePostMessage(id int, newPost models.Post) (models.Post, error)
}
