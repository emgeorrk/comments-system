package storage

import (
	"github.com/google/uuid"
	"graphql-comments/types"
)

type DataStore interface {
	AddPost(title, content string) *types.Post
	AddComment(postID, parentCommentID, content string) (*types.Comment, error)
	GetPosts() []*types.Post
	GetPostByID(id string) (*types.Post, error)
	GetCommentByID(id string) (*types.Comment, error)
	GetComments(postID string, paginationSize int) ([]*types.Comment, error)
}

var DataBase DataStore

func GenerateNewPostUUID() string {
	for {
		newUUID := uuid.New()
		if _, err := DataBase.GetPostByID(newUUID.String()); err != nil {
			return newUUID.String()
		}
	}
}

func GenerateNewCommentUUID() string {
	for {
		newUUID := uuid.New()
		if _, err := DataBase.GetCommentByID(newUUID.String()); err != nil {
			return newUUID.String()
		}
	}
}
