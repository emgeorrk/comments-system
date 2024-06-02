package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"graphql-comments/storage"
	"graphql-comments/types"
	"time"
)

type DataStorePostgres struct {
	DB *sql.DB
}

func NewPostgresDataStore(dbURL string) (*DataStorePostgres, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return &DataStorePostgres{DB: db}, nil
}

func (store *DataStorePostgres) AddPost(title, content string, allowComments bool) (*types.Post, error) {
	if title == "" {
		return nil, errors.New("title is empty")
	}
	if content == "" {
		return nil, errors.New("content is empty")
	}

	post := &types.Post{
		ID:            storage.GenerateNewPostUUID(),
		Title:         title,
		Content:       content,
		CreatedAt:     time.Now(),
		Comments:      []string{},
		AllowComments: allowComments,
	}

	_, err := store.DB.Exec("INSERT INTO posts (id, title, content, created_at) VALUES ($1, $2, $3, $4)",
		post.ID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (store *DataStorePostgres) AddComment(postID, parentCommentID, content string) (*types.Comment, error) {
	switch {
	case postID == "":
		return nil, errors.New("postID is empty")
	case content == "":
		return nil, errors.New("content is empty")
	case len(content) > storage.MaxCommentLength:
		return nil, errors.New(fmt.Sprintf("content is too long (maximum %d chars)", storage.MaxCommentLength))
	}

	post, err := store.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	if !post.AllowComments {
		return nil, errors.New("comments are not allowed for this post")
	}

	comment := &types.Comment{
		ID:              storage.GenerateNewCommentUUID(),
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Content:         content,
		CreatedAt:       time.Now(),
		Replies:         []string{},
	}

	if _, err := store.DB.Exec(
		"INSERT INTO comments (id, post_id, parent_comment_id, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		comment.ID, comment.PostID, comment.ParentCommentID, comment.Content, comment.CreatedAt,
	); err != nil {
		return nil, err
	}

	if parentCommentID == "" {
		// Добавление комментария к посту
		_, err = store.DB.Exec("UPDATE posts SET comments = array_append(comments, $1) WHERE id = $2", comment.ID, postID)
		if err != nil {
			return nil, err
		}
	} else {
		// Добавление вложенного комментария
		_, err = store.DB.Exec(
			"UPDATE comments SET replies = array_append(replies, $1) WHERE id = $2", comment.ID, parentCommentID)
		if err != nil {
			return nil, err
		}
	}

	return comment, nil
}

func (store *DataStorePostgres) GetPosts() ([]*types.Post, error) {
	rows, err := store.DB.Query("SELECT id, title, content, created_at, comments, allow_comments FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*types.Post, 0)
	for rows.Next() {
		post := &types.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Comments, &post.AllowComments)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (store *DataStorePostgres) GetPostByID(id string) (*types.Post, error) {
	post := &types.Post{}
	err := store.DB.QueryRow("SELECT id, title, content, created_at, comments, allow_comments FROM posts WHERE id = $1", id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Comments,
		&post.AllowComments,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return post, nil
}

func (store *DataStorePostgres) GetComments(postID string) ([]*types.Comment, error) {
	rows, err := store.DB.Query("SELECT id, post_id, parent_comment_id, content, created_at FROM comments WHERE post_id = $1", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*types.Comment, 0)
	for rows.Next() {
		comment := &types.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentCommentID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (store *DataStorePostgres) GetCommentByID(id string) (*types.Comment, error) {
	comment := &types.Comment{}
	err := store.DB.QueryRow(
		"SELECT id, post_id, parent_comment_id, content, created_at FROM comments WHERE id = $1", id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.ParentCommentID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return comment, nil
}

func (store *DataStorePostgres) GetReplies(commentID string) ([]*types.Comment, error) {
	rows, err := store.DB.Query("SELECT id, post_id, parent_comment_id, content, created_at FROM comments WHERE parent_comment_id = $1", commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*types.Comment, 0)
	for rows.Next() {
		comment := &types.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentCommentID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
