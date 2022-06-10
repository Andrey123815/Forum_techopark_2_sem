package forumRepo

import (
	"db-forum/internal/app/models"
	"fmt"
	"github.com/jackc/pgx"
	"time"
)

type Repository struct {
	Database *pgx.ConnPool
}

const EMPTY_PARAMETER = -1

func CreateForumRepository(db *pgx.ConnPool) *Repository {
	return &Repository{Database: db}
}

func (repository *Repository) GetForumBySlug(slug string) (models.Forum, error) {
	var forum models.Forum

	err := repository.Database.QueryRow(`SELECT id,title,"user",slug,posts,threads FROM forums WHERE slug=$1;`, slug).
		Scan(&forum.Id, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}

func (repository *Repository) CreateNewForum(title, user, slug string) (models.Forum, error) {
	var newForum models.Forum

	err := repository.Database.QueryRow(`INSERT INTO forums(title, "user", slug) VALUES ($1, $2, $3)
		RETURNING title, "user", slug, posts, threads;`, title, user, slug).
		Scan(&newForum.Title, &newForum.User, &newForum.Slug, &newForum.Posts, &newForum.Threads)
	if err != nil {
		return models.Forum{}, err
	}

	return newForum, nil
}

func (repository *Repository) CreateNewThread(title, author, forum, message, slug string, created time.Time) (models.Thread, error) {
	var newThread models.Thread

	err := repository.Database.QueryRow(`INSERT INTO threads(title,author,forum,message,slug,created) VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, title, author, forum, message, votes, slug, created;`, title, author, forum, message, slug, created).
		Scan(&newThread.Id, &newThread.Title, &newThread.Author, &newThread.Forum,
			&newThread.Message, &newThread.Votes, &newThread.Slug, &newThread.Created)
	if err != nil {
		return models.Thread{}, err
	}

	return newThread, nil
}

func (repository *Repository) GetAlreadyExistThread(slug string) (models.Thread, error) {
	var existThread models.Thread

	if slug == "" {
		return models.Thread{}, nil
	}

	err := repository.Database.QueryRow(`SELECT * FROM threads WHERE slug = $1;`, slug).
		Scan(&existThread.Id, &existThread.Title, &existThread.Author, &existThread.Forum,
			&existThread.Message, &existThread.Votes, &existThread.Slug, &existThread.Created)
	if err != nil {
		return models.Thread{}, err
	}

	return existThread, nil
}

func (repository *Repository) GetThreadsBySlug(forumSlug string, limit int, desc bool, since string) ([]models.Thread, error) {
	threads := make([]models.Thread, 0, 0)

	query := `SELECT * FROM threads WHERE forum = $1`

	if since != "" {
		comparator := ">="
		if desc == true {
			comparator = "<="
		}
		query += fmt.Sprintf(` AND created %s '%s'`, comparator, since)
	}

	sortDirection := "ASC"
	if desc == true {
		sortDirection = "DESC"
	}

	if limit > 0 {
		query += fmt.Sprintf(` ORDER BY created %s LIMIT %d;`, sortDirection, limit)
	} else {
		query += fmt.Sprintf(` ORDER BY created %s;`, sortDirection)
	}

	result, err := repository.Database.Query(query, forumSlug)
	if err != nil {
		fmt.Println(err)
		return threads, err
	}

	defer func() {
		if result != nil {
			result.Close()
		}
	}()

	for result.Next() {
		var thread models.Thread
		err := result.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
		if err != nil {
			return []models.Thread{}, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (repository *Repository) GetForumUsers(forumID int64, limit int, sortDirection string, since string) ([]models.User, error) {
	users := make([]models.User, 0)

	query := `SELECT "nickname", "about", "email", "fullname" FROM "users"
				WHERE "id" IN (SELECT "nickname" FROM "forum_users" WHERE forum = $1)`
	if since != "" {
		sign := ">"
		if sortDirection == "DESC" {
			sign = "<"
		}
		query += fmt.Sprintf(` AND "nickname" %s '%s'`, sign, since)
	}

	if limit == EMPTY_PARAMETER {
		limit = 1000
	}
	query += fmt.Sprintf(` ORDER BY "nickname" %s LIMIT %d;`, sortDirection, limit)

	result, err := repository.Database.Query(query, forumID)

	if err != nil {
		return []models.User{}, err
	}
	defer result.Close()

	for result.Next() {
		var user models.User
		err := result.Scan(
			&user.Nickname,
			&user.About,
			&user.Email,
			&user.Fullname,
		)
		if err != nil {
			return []models.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}
