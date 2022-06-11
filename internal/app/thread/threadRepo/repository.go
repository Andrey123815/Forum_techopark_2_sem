package threadRepo

import (
	"db-forum/internal/app/models"
	"db-forum/internal/customErrors"
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
	"time"
)

type Repository struct {
	Database *pgx.ConnPool
}

func CreateThreadRepository(db *pgx.ConnPool) *Repository {
	return &Repository{Database: db}
}

const POST_PARAMS = `id, parent, author, message, "isEdited", forum, thread, created`

func generateFlatSortQuery(limit, since, sortDirection string) string {
	query := fmt.Sprintf(`SELECT %s FROM posts WHERE thread = $1`, POST_PARAMS)

	comp := "<"
	if sortDirection == "ASC" {
		comp = ">"
	}

	if since != "" {
		query += fmt.Sprintf(` AND id %s %s`, comp, since)
	}
	query += fmt.Sprintf(" ORDER BY created %s, id %s LIMIT %s;", sortDirection, sortDirection, limit)

	return query
}

func generateTreeSortQuery(limit, since, sortDirection string) string {
	query := fmt.Sprintf(`SELECT %s FROM posts WHERE thread = $1`, POST_PARAMS)

	comp := "<"
	if sortDirection == "ASC" {
		comp = ">"
	}

	if since != "" {
		query += fmt.Sprintf(` AND path %s (SELECT path FROM posts WHERE id = %s) `, comp, since)
	}
	query += fmt.Sprintf(` ORDER BY path[1] %s, path %s LIMIT %s;`, sortDirection, sortDirection, limit)

	return query
}

func generateParentTreeSortQuery(limit string, since string, sortDirection string) string {
	query := fmt.Sprintf(`SELECT %s FROM posts WHERE thread = $1 AND path &&
				(SELECT ARRAY (SELECT id FROM posts WHERE thread = $1 AND parent = 0 `, POST_PARAMS)

	comp := "<"
	if sortDirection == "ASC" {
		comp = ">"
	}

	if since != "" {
		query += fmt.Sprintf(`AND path %s (SELECT path[1:1] FROM posts WHERE id = %s LIMIT 1) `, comp, since)
	}
	query += fmt.Sprintf(`ORDER BY path[1] %s, path LIMIT %s)) ORDER BY path[1] %s, path;`, sortDirection, limit, sortDirection)

	return query
}

func (repository *Repository) GetThreadByID(threadID int64) (models.Thread, error) {
	var thread models.Thread

	err := repository.Database.QueryRow(`SELECT * FROM threads WHERE id = $1;`, threadID).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
			&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func (repository *Repository) GetThreadBySlugOrID(slugOrID string) (models.Thread, error) {
	var thread models.Thread

	threadID, _ := strconv.Atoi(slugOrID)

	err := repository.Database.QueryRow(`SELECT * FROM threads WHERE id = $1 OR slug = $2;`, threadID, slugOrID).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
			&thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func (repository *Repository) GetThreadIDBySlugOrID(slugOrID string) (int64, error) {
	var id int64

	threadID, _ := strconv.Atoi(slugOrID)

	err := repository.Database.QueryRow(`SELECT id FROM threads WHERE id = $1 OR slug = $2;`, threadID, slugOrID).
		Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (repository *Repository) CreateNewPosts(newPosts []models.Post, thread models.Thread) ([]models.Post, error) {
	insertQuery := "($%d, $%d, $%d, $%d, $%d, $%d)"
	query := `INSERT INTO posts (parent, author, message, forum, thread, created) VALUES `

	var insertionData []interface{}
	dataCreation := time.Now()
	row, rowSize := 0, 6

	for i, singlePost := range newPosts {
		err := repository.Database.QueryRow(`SELECT nickname FROM users WHERE nickname = $1;`, singlePost.Author).Scan(&singlePost.Author)
		if err != nil {
			return []models.Post{}, customErrors.NoAuthor
		}

		if singlePost.Parent != 0 {
			var id int64
			err = repository.Database.QueryRow(`SELECT id FROM posts WHERE thread = $1 AND id = $2;`, thread.Id, singlePost.Parent).Scan(&id)
			if err != nil {
				return []models.Post{}, customErrors.AnotherThreadPost
			}
		}

		if i != 0 {
			query += ", "
		}

		query += fmt.Sprintf(insertQuery, row+1, row+2, row+3, row+4, row+5, row+6)

		row += rowSize

		insertionData = append(insertionData, singlePost.Parent, singlePost.Author, singlePost.Message, thread.Forum, thread.Id, dataCreation)
	}

	query += `RETURNING id, parent, author, message, "isEdited", forum, thread, created;`

	result, err := repository.Database.Query(query, insertionData...)
	if err != nil {
		return []models.Post{}, err
	}

	defer func() {
		if result != nil {
			result.Close()
		}
	}()

	var results []models.Post

	for result.Next() {
		var singlePost models.Post
		err = result.Scan(&singlePost.Id, &singlePost.Parent, &singlePost.Author,
			&singlePost.Message, &singlePost.IsEdited, &singlePost.Forum,
			&singlePost.Thread, &singlePost.Created,
		)
		if err != nil {
			return []models.Post{}, err
		}

		results = append(results, singlePost)
	}

	return results, nil
}

func (repository *Repository) VoteThread(nickname string, threadID int64, voice int32, oldTread *models.Thread) error {
	var (
		vote     models.Vote
		user     string
		thread   int64
		oldVoice int32
	)

	err := repository.Database.QueryRow(`SELECT "user", thread, voice FROM votes WHERE "user" = $1 AND thread = $2;`,
		nickname, threadID).Scan(&user, &thread, &oldVoice)

	if err != nil {
		oldTread.Votes += voice
		err := repository.Database.QueryRow(`INSERT INTO votes ("user", thread, voice) VALUES ($1, $2, $3)
		RETURNING "user", voice;`, nickname, threadID, voice).
			Scan(&vote.Nickname, &vote.Voice)
		if err != nil {
			return err
		}

		err = repository.Database.QueryRow(`UPDATE threads SET votes = (votes+$1) WHERE id = $2 RETURNING votes;`,
			voice, threadID).Scan(&thread)

		return nil
	}

	if voice == oldVoice {
		return nil
	}

	oldTread.Votes += 2 * voice
	err = repository.Database.QueryRow(`UPDATE threads SET votes = votes + $1 WHERE id = $2 RETURNING votes;`,
		2*voice, threadID).Scan(&thread)

	err = repository.Database.QueryRow(`UPDATE votes SET voice=$1 WHERE "user" = $2 AND thread = $3 RETURNING "user";`,
		voice, nickname, threadID).Scan(&user)

	return err
}

func (repository *Repository) GetThreadPosts(id int64, limit, since, sortType, sortDirection string) ([]models.Post, error) {
	posts := make([]models.Post, 0)

	var query string

	switch sortType {
	case "flat":
		query = generateFlatSortQuery(limit, since, sortDirection)
	case "tree":
		query = generateTreeSortQuery(limit, since, sortDirection)
	case "parent_tree":
		query = generateParentTreeSortQuery(limit, since, sortDirection)
	}

	result, err := repository.Database.Query(query, id)
	if err != nil {
		return []models.Post{}, err
	}
	defer result.Close()

	for result.Next() {
		var post models.Post
		err := result.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return []models.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository *Repository) UpdateThreadDetails(threadID int64, title, message string) (models.Thread, error) {
	var updatedThread models.Thread

	err := repository.Database.QueryRow(`UPDATE threads SET title = $1, message = $2 WHERE id = $3
											RETURNING id, title, author, forum, message, votes, slug, created;`,
		title, message, threadID).Scan(&updatedThread.Id, &updatedThread.Title, &updatedThread.Author,
		&updatedThread.Forum, &updatedThread.Message, &updatedThread.Votes, &updatedThread.Slug, &updatedThread.Created)

	return updatedThread, err
}
