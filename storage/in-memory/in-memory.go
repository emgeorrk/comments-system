package inMemory

import (
	"errors"
	"fmt"
	"graphql-comments/storage"
	"graphql-comments/types"
	"time"
)

// InMemoryStore структура для хранения постов и комментариев в памяти
type InMemoryStore struct {
	Posts    map[string]*types.Post
	Comments map[string]*types.Comment
}

// NewInMemoryStore создает новый in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Posts:    make(map[string]*types.Post),
		Comments: make(map[string]*types.Comment),
	}
}

func (store *InMemoryStore) AddPost(title, content string) (*types.Post, error) {
	if title == "" {
		return nil, errors.New("title is empty")
	}
	if content == "" {
		return nil, errors.New("content is empty")
	}

	post := &types.Post{
		ID:        storage.GenerateNewPostUUID(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		Comments:  []string{},
	}
	store.Posts[post.ID] = post
	return post, nil
}

func (store *InMemoryStore) AddComment(postID, parentCommentID string, content string) (*types.Comment, error) {
	if content == "" {
		return nil, errors.New("content is empty")
	}
	if len(content) > storage.MaxCommentLength {
		return nil, errors.New(fmt.Sprintf("content is too long (maximum %d chars)", storage.MaxCommentLength))
	}

	comment := &types.Comment{
		ID:              storage.GenerateNewCommentUUID(),
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Content:         content,
		CreatedAt:       time.Now(),
		Replies:         []string{},
	}

	store.Comments[comment.ID] = comment

	if parentCommentID == "" {
		// Добавление комментария к посту
		if post, ok := store.Posts[postID]; ok {
			post.Comments = append(post.Comments, comment.ID)
		} else {
			return nil, errors.New("post not found")
		}
	} else {
		// Добавление вложенного комментария
		if parentComment, ok := store.Comments[parentCommentID]; ok {
			parentComment.Replies = append(parentComment.Replies, comment.ID)
		} else {
			return nil, errors.New("parent comment not found")
		}
	}
	return comment, nil
}

func (store *InMemoryStore) GetPosts() ([]*types.Post, error) {
	posts := make([]*types.Post, 0)

	for _, post := range store.Posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (store *InMemoryStore) GetPostByID(id string) (*types.Post, error) {
	if post, ok := store.Posts[id]; ok {
		return post, nil
	}
	return nil, errors.New("post not found")
}

func (store *InMemoryStore) GetComments(postID string) ([]*types.Comment, error) {
	post, err := store.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	comments := make([]*types.Comment, 0)
	for _, commentID := range post.Comments {
		comment, err := store.GetCommentByID(commentID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (store *InMemoryStore) GetCommentByID(id string) (*types.Comment, error) {
	if comment, ok := store.Comments[id]; ok {
		return comment, nil
	}
	return nil, errors.New("comment not found")
}

func (store *InMemoryStore) GetReplies(commentID string) ([]*types.Comment, error) {
	comment, err := store.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	replies := make([]*types.Comment, 0)
	for _, replyID := range comment.Replies {
		reply, err := store.GetCommentByID(replyID)
		if err != nil {
			return nil, err
		}

		replies = append(replies, reply)
	}

	return replies, nil
}
