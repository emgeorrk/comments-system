package postgres

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"graphql-comments/storage"
	"graphql-comments/types"
	"time"
)

type PostgresDataStore struct {
	db *sql.DB
}

func NewPostgresDataStore(dbURL string) (*PostgresDataStore, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return &PostgresDataStore{db: db}, nil
}

func (store *PostgresDataStore) AddPost(title, content string) (*types.Post, error) {
	post := &types.Post{
		ID:        storage.GenerateNewPostUUID(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		Comments:  []string{},
	}
	
	_, err := store.db.Exec("INSERT INTO posts (id, title, content, created_at) VALUES ($1, $2, $3, $4)",
		post.ID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return post, nil
}

func (store *PostgresDataStore) AddComment(postID, parentCommentID, content string) (*types.Comment, error) {
	return nil, nil
}

func (store *PostgresDataStore) GetPosts() ([]*types.Post, error) {
	rows, err := store.db.Query("SELECT id, title, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	posts := make([]*types.Post, 0)
	for rows.Next() {
		post := &types.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
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

func (store *PostgresDataStore) GetPostByID(id string) (*types.Post, error) {
	post := &types.Post{}
	err := store.db.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = $1", id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return post, nil
}

func (store *PostgresDataStore) GetComments(postID string) ([]*types.Comment, error) {
	rows, err := store.db.Query("SELECT id, post_id, parent_comment_id, content, created_at FROM comments WHERE post_id = $1", postID)
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

func (store *PostgresDataStore) GetCommentByID(id string) (*types.Comment, error) {
	comment := &types.Comment{}
	err := store.db.QueryRow(
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

func (store *PostgresDataStore) GetReplies(commentID string) ([]*types.Comment, error) {
	return nil, nil
}
