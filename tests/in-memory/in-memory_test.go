package inMemory_test

import (
	"graphql-comments/storage"
	inMemory "graphql-comments/storage/in-memory"
	"testing"
	_ "time"
)

func TestAddPost(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	t.Run("AddPost", func(t *testing.T) {
		post, err := store.AddPost("Test Title", "Test Content")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if post.Title != "Test Title" || post.Content != "Test Content" {
			t.Errorf("Post fields do not match the input")
		}

		if _, ok := store.Posts[post.ID]; !ok {
			t.Errorf("Post was not added to the store")
		}
	})

	t.Run("AddPostWithEmptyTitle", func(t *testing.T) {
		_, err := store.AddPost("", "Test Content")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("AddPostWithEmptyContent", func(t *testing.T) {
		_, err := store.AddPost("Test Title", "")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetPostByID(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	post, _ := store.AddPost("Test Title", "Test Content")

	t.Run("GetPostByIDSuccessfully", func(t *testing.T) {
		retrievedPost, err := store.GetPostByID(post.ID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if retrievedPost.ID != post.ID {
			t.Errorf("Retrieved post ID does not match the input")
		}
	})

	t.Run("GetPostByIDWithNonexistentID", func(t *testing.T) {
		_, err := store.GetPostByID("nonexistent-id")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestAddComment(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	post, _ := store.AddPost("Test Title", "Test Content")

	t.Run("AddCommentSuccessfully", func(t *testing.T) {
		comment, err := store.AddComment(post.ID, "", "Test Comment")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if comment.Content != "Test Comment" {
			t.Errorf("Comment content does not match the input")
		}

		if _, ok := store.Comments[comment.ID]; !ok {
			t.Errorf("Comment was not added to the store")
		}
	})

	t.Run("AddCommentWithNonexistentPostID", func(t *testing.T) {
		_, err := store.AddComment("nonexistent-id", "", "Test Comment")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("AddCommentWithEmptyContent", func(t *testing.T) {
		_, err := store.AddComment(post.ID, "", "")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
