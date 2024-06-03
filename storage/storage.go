package storage

import (
	"github.com/google/uuid"
	"graphql-comments/types"
)

const (
	MaxCommentLength     = 2000
	MaxPostTitleLength   = 100
	MaxPostContentLength = 10000
	CommentsPageSize     = 10
)

type DataStore interface {
	AddPost(title, content string, allowComments bool) (*types.Post, error)
	AddComment(postID, parentCommentID, content string) (*types.Comment, error)
	GetPosts() ([]*types.Post, error)
	GetPostByID(id string) (*types.Post, error)
	GetComments(postID string, page int) ([]*types.Comment, error)
	GetCommentByID(id string) (*types.Comment, error)
	GetNumberOfCommentPages(postID string) (int, error)
	GetReplies(commentID string) ([]*types.Comment, error)
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
