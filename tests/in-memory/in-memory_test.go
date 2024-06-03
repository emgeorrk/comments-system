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
		post, err := store.AddPost("Test Title", "Test Content", true)
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
}

func TestGetPostByID(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	post, _ := store.AddPost("Test Title", "Test Content", true)

	t.Run("GetPostByID", func(t *testing.T) {
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

	post, _ := store.AddPost("Test Title", "Test Content", true)

	t.Run("AddComment", func(t *testing.T) {
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
}

func TestGetPosts(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	t.Run("GetPostsReturnsEmptySliceWhenNoPosts", func(t *testing.T) {
		posts, err := store.GetPosts()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(posts) != 0 {
			t.Errorf("Expected 0 posts, got %v", len(posts))
		}
	})

	t.Run("GetPosts", func(t *testing.T) {
		post1, _ := store.AddPost("Title 1", "Content 1", true)
		post2, _ := store.AddPost("Title 2", "Content 2", true)

		posts, err := store.GetPosts()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(posts) != 2 {
			t.Errorf("Expected 2 posts, got %v", len(posts))
		}

		if posts[0] != post1 || posts[1] != post2 {
			t.Errorf("Posts returned are not the same as the ones added")
		}
	})
}

func TestGetComments(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	t.Run("GetComments", func(t *testing.T) {
		post, _ := store.AddPost("Title", "Content", true)

		_, err := store.AddComment(post.ID, "", "Comment 1")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		_, err = store.AddComment(post.ID, "", "Comment 2")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		comments, err := store.GetComments(post.ID, 1)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(comments) != 2 {
			t.Errorf("Expected 2 comments, got %v", len(comments))
		}
	})

	t.Run("GetCommentsWhenNoComments", func(t *testing.T) {
		post, _ := store.AddPost("Title", "Content", true)
		store.AddComment(post.ID, "", "Comment 1")
		store.AddComment(post.ID, "", "Comment 2")

		comments, err := store.GetComments(post.ID, 1)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(comments) != 2 {
			t.Errorf("Expected 2 comments, got %v", len(comments))
		}
	})

	t.Run("GetCommentsWithNonexistentPostID", func(t *testing.T) {
		_, err := store.GetComments("nonexistent-id", 1)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetCommentByID(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	t.Run("GetCommentByID", func(t *testing.T) {
		post, _ := store.AddPost("Title", "Content", true)
		comment, _ := store.AddComment(post.ID, "", "Comment")

		retrievedComment, err := store.GetCommentByID(comment.ID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if retrievedComment.ID != comment.ID {
			t.Errorf("Retrieved comment ID does not match the input")
		}
	})

	t.Run("GetCommentByIDWithNonexistentID", func(t *testing.T) {
		_, err := store.GetCommentByID("nonexistent-id")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetNumberOfCommentPages(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	t.Run("GetNumberOfCommentPages", func(t *testing.T) {
		post, _ := store.AddPost("Title", "Content", true)

		for i := 0; i < storage.CommentsPageSize*3; i++ {
			_, err := store.AddComment(post.ID, "", "Comment")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}

		numPages, err := store.GetNumberOfCommentPages(post.ID)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if numPages != 3 {
			t.Errorf("Expected 3 pages, got %v", numPages)
		}
	})

	t.Run("GetNumberOfCommentPagesWithNonexistentPostID", func(t *testing.T) {
		_, err := store.GetNumberOfCommentPages("nonexistent-id")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetReplies(t *testing.T) {
	store := inMemory.NewInMemoryStore()
	storage.DataBase = store

	t.Run("GetReplies", func(t *testing.T) {
		post, _ := store.AddPost("Title", "Content", true)
		comment1, _ := store.AddComment(post.ID, "", "Comment 1")
		comment2, _ := store.AddComment(post.ID, comment1.ID, "Comment 2")

		replies, err := store.GetReplies(comment1.ID)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(replies) != 1 {
			t.Errorf("Expected 1 reply, got %v", len(replies))
		}

		if replies[0].ID != comment2.ID {
			t.Errorf("Retrieved reply ID does not match the input")
		}
	})

	t.Run("GetRepliesWithNonexistentCommentID", func(t *testing.T) {
		_, err := store.GetReplies("nonexistent-id")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
