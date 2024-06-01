package types

import "time"

// Post структура для хранения постов
type Post struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	Comments  []*Comment
}

// Comment структура для хранения комментариев
type Comment struct {
	ID              string
	PostID          string
	ParentCommentID string
	Content         string
	CreatedAt       time.Time
	Replies         []*Comment
}
