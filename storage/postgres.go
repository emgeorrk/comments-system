package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"graphql-comments/types"
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

func (store *PostgresDataStore) AddPost(title, content string) *types.Post {
	return nil
}

func (store *PostgresDataStore) AddComment(postID, parentCommentID, content string) (*types.Comment, error) {
	return nil, nil
}

func (store *PostgresDataStore) GetPostByID(id string) (*types.Post, error) {
	return nil, nil
}

func (store *PostgresDataStore) GetPosts() []*types.Post {
	return nil
}

func (store *PostgresDataStore) GetCommentByID(id string) (*types.Comment, error) {
	return nil, nil
}

func (store *PostgresDataStore) GetComments(postID string) ([]*types.Comment, error) {
	return nil, nil
}

func (store *PostgresDataStore) GetReplies(commentID string) ([]*types.Comment, error) {
	return nil, nil
}
