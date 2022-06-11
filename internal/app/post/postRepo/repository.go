package postRepo

import (
	"db-forum/internal/app/models"
	"github.com/jackc/pgx"
)

type Repository struct {
	Database *pgx.ConnPool
}

func CreatePostRepository(db *pgx.ConnPool) *Repository {
	return &Repository{Database: db}
}

func (r *Repository) GetPostDetails(id int, related []string) (postInfo models.DetailPostInfo, err error) {
	var post models.Post
	err = r.Database.QueryRow(`SELECT id, parent, author, message, "isEdited", forum, thread, created FROM posts WHERE id = $1;`, id).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return
	}

	postInfo.Post = &post

	if len(related) != 0 {
		for _, filter := range related {
			switch filter {
			case "user":
				var user models.User
				err = r.Database.QueryRow(`SELECT nickname, fullname, about, email FROM users WHERE nickname = $1;`, post.Author).
					Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
				if err != nil {
					return
				}
				postInfo.Author = &user
			case "forum":
				var forum models.Forum
				err = r.Database.QueryRow(`SELECT title, "user", slug, posts, threads FROM forums WHERE slug = $1;`, post.Forum).
					Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
				if err != nil {
					return
				}
				postInfo.Forum = &forum
			case "thread":
				var thread models.Thread
				err = r.Database.QueryRow(`SELECT id, title, author,forum, message, votes, slug,
						created FROM threads WHERE id = $1;`, post.Thread).
					Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
						&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
				if err != nil {
					return
				}
				postInfo.Thread = &thread
			}
		}
	}

	return
}

func (r *Repository) ChangePostMessage(id int, newPost models.Post) (models.Post, error) {
	var post models.Post
	err := r.Database.QueryRow(`UPDATE posts SET message = $1, "isEdited" = true WHERE id = $2
			RETURNING id, parent, author, message, "isEdited", forum, thread, created;`, newPost.Message, id).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	return post, err
}
