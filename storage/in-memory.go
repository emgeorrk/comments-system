package storage

import (
	"errors"
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

// AddPost добавляет новый пост
func (store *InMemoryStore) AddPost(title, content string) *types.Post {
	post := &types.Post{
		ID:        GenerateNewPostUUID(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		Comments:  []*types.Comment{},
	}
	store.Posts[post.ID] = post
	return post
}

// AddComment добавляет новый комментарий
func (store *InMemoryStore) AddComment(postID, parentCommentID string, content string) (*types.Comment, error) {
	comment := &types.Comment{
		ID:              GenerateNewCommentUUID(),
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Content:         content,
		CreatedAt:       time.Now(),
		Replies:         []*types.Comment{},
	}
	
	store.Comments[comment.ID] = comment
	
	if parentCommentID == "" {
		// Добавление комментария к посту
		if post, ok := store.Posts[postID]; ok {
			post.Comments = append(post.Comments, comment)
		} else {
			return nil, errors.New("post not found")
		}
	} else {
		// Добавление вложенного комментария
		if parentComment, ok := store.Comments[parentCommentID]; ok {
			parentComment.Replies = append(parentComment.Replies, comment)
		} else {
			return nil, errors.New("parent comment not found")
		}
	}
	return comment, nil
}

// GetPostByID возвращает пост по ID
func (store *InMemoryStore) GetPostByID(id string) (*types.Post, error) {
	if post, ok := store.Posts[id]; ok {
		return post, nil
	}
	return nil, errors.New("post not found")
}

func (store *InMemoryStore) GetCommentByID(id string) (*types.Comment, error) {
	if comment, ok := store.Comments[id]; ok {
		return comment, nil
	}
	return nil, errors.New("comment not found")
}

// GetComments получает комментарии к посту
func (store *InMemoryStore) GetComments(postID string, paginationSize int) ([]*types.Comment, error) {
	if post, ok := store.Posts[postID]; ok {
		return post.Comments, nil
	}
	return nil, errors.New("post not found")
}

// GetPosts возвращает все существующие посты
func (store *InMemoryStore) GetPosts() []*types.Post {
	posts := make([]*types.Post, 0)
	
	for _, post := range store.Posts {
		posts = append(posts, post)
	}
	
	return posts
}
