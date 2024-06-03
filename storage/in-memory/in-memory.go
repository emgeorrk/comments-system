package inMemory

import (
	"errors"
	"graphql-comments/storage"
	"graphql-comments/types"
	"time"
)

// DataStoreInMemory структура для хранения постов и комментариев в памяти
type DataStoreInMemory struct {
	Posts    map[string]*types.Post
	Comments map[string]*types.Comment
}

// NewInMemoryStore создает новый in-memory store
func NewInMemoryStore() *DataStoreInMemory {
	return &DataStoreInMemory{
		Posts:    make(map[string]*types.Post),
		Comments: make(map[string]*types.Comment),
	}
}

func (store *DataStoreInMemory) AddPost(title, content string, allowComments bool) (*types.Post, error) {
	post := &types.Post{
		ID:            storage.GenerateNewPostUUID(),
		Title:         title,
		Content:       content,
		CreatedAt:     time.Now(),
		Comments:      []string{},
		AllowComments: allowComments,
	}
	store.Posts[post.ID] = post
	return post, nil
}

func (store *DataStoreInMemory) AddComment(postID, parentCommentID string, content string) (*types.Comment, error) {
	post, ok := store.Posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}

	if !store.Posts[postID].AllowComments {
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

	store.Comments[comment.ID] = comment

	if parentCommentID == "" {
		// Добавление комментария к посту
		post.Comments = append(post.Comments, comment.ID)
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

func (store *DataStoreInMemory) GetPosts() ([]*types.Post, error) {
	posts := make([]*types.Post, 0)

	for _, post := range store.Posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (store *DataStoreInMemory) GetPostByID(id string) (*types.Post, error) {
	if post, ok := store.Posts[id]; ok {
		return post, nil
	}
	return nil, errors.New("post not found")
}

func (store *DataStoreInMemory) GetComments(postID string, page int) ([]*types.Comment, error) {
	post, err := store.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	comments := make([]*types.Comment, 0)
	for idx := (page - 1) * storage.CommentsPageSize; idx < len(post.Comments) && idx < page*storage.CommentsPageSize; idx++ {
		comment, err := store.GetCommentByID(post.Comments[idx])
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (store *DataStoreInMemory) GetCommentByID(id string) (*types.Comment, error) {
	if comment, ok := store.Comments[id]; ok {
		return comment, nil
	}
	return nil, errors.New("comment not found")
}

func (store *DataStoreInMemory) GetNumberOfCommentPages(postID string) (int, error) {
	if _, ok := store.Posts[postID]; !ok {
		return 0, errors.New("post not found")
	}
	return len(store.Posts[postID].Comments) / storage.CommentsPageSize, nil
}

func (store *DataStoreInMemory) GetReplies(commentID string) ([]*types.Comment, error) {
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
