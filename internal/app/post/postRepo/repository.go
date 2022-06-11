package postRepo

import (
	"db-forum/internal/app/models"
	"github.com/jackc/pgx"
)

type Repository struct {
	Database *pgx.ConnPool
}

const EMPTY_PARAMETER = -1

func CreatePostRepository(db *pgx.ConnPool) *Repository {
	return &Repository{Database: db}
}

func (repository *Repository) GetPostDetails(postID int64, relatedArr []string) (map[string]interface{}, error) {
	var postInfo models.Post

	err := repository.Database.QueryRow(`SELECT id,author,message,"isEdited",forum,thread,created FROM posts WHERE id = $1;`, postID).
		Scan(&postInfo.Id, &postInfo.Author, &postInfo.Message, &postInfo.IsEdited, &postInfo.Forum, &postInfo.Thread, &postInfo.Created)
	if err != nil {
		return map[string]interface{}{}, err
	}

	postDetails := map[string]interface{}{
		"post": postInfo,
	}

	if len(relatedArr) != 0 {
		for _, related := range relatedArr {

			switch related {
			case "user":
				var user models.User
				err = repository.Database.QueryRow(`SELECT nickname,fullname,about,email FROM users WHERE nickname = $1;`, postInfo.Author).
					Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
				if err != nil {
					return postDetails, err
				}
				postDetails["author"] = user

			case "forum":
				var forum models.Forum
				err = repository.Database.QueryRow(`SELECT id,title,"user",slug,posts,threads
								FROM forums WHERE slug = (SELECT forum FROM posts WHERE id = $1);`, postInfo.Id).
					Scan(&forum.Id, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
				if err != nil {
					return postDetails, err
				}
				postDetails["forum"] = forum

			case "thread":
				var thread models.Thread
				err = repository.Database.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created
								FROM threads WHERE id = (SELECT thread FROM posts WHERE id = $1);`, postInfo.Id).
					Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
						&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
				if err != nil {
					return postDetails, err
				}
				postDetails["thread"] = thread
			}
		}
	}

	return postDetails, nil
}

func (repository *Repository) ChangePostMessage(postID int64, message string) (models.Post, error) {
	var updatedPost models.Post

	err := repository.Database.QueryRow(`UPDATE posts SET message = $1, "isEdited" = $2
			WHERE id = $3 RETURNING id,author,message,"isEdited",forum,thread,created;`, message, true, postID).
		Scan(&updatedPost.Id, &updatedPost.Author, &updatedPost.Message, &updatedPost.IsEdited,
			&updatedPost.Forum, &updatedPost.Thread, &updatedPost.Created)
	if err != nil {
		return models.Post{}, err
	}

	return updatedPost, nil
}
