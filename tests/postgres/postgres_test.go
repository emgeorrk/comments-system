package postgres_test

import (
	"database/sql"
	"errors"
	"fmt"
	"graphql-comments/storage"
	"graphql-comments/storage/postgres"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func TestAddPost(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	t.Run("AddPost", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO posts (id, title, content, created_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(sqlmock.AnyArg(), "Test Title", "Test Content", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		_, err := store.AddPost("Test Title", "Test Content")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("AddPostWithEmptyTitle", func(t *testing.T) {
		title := ""
		content := "Test Content"

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO posts (id, title, content, created_at) VALUES ($1, $2, $3, $4)")).
			WithArgs(sqlmock.AnyArg(), title, content, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := store.AddPost(title, content)
		fmt.Println(err)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestAddComment(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	t.Run("AddCommentToPost", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO comments")).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta("UPDATE posts SET comments")).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err := store.AddComment("post-id", "", "Test Comment")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestAddCommentWithDBError(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	mock.ExpectExec("INSERT INTO comments").WillReturnError(errors.New("DB error"))

	_, err := store.AddComment("post-id", "", "Test Comment")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetPosts(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow("post-id", "Test Title", "Test Content", time.Now())

	mock.ExpectQuery("SELECT id, title, content, created_at FROM posts").WillReturnRows(rows)

	_, err := store.GetPosts()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetPostsWithDBError(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	mock.ExpectQuery("SELECT id, title, content, created_at FROM posts").WillReturnError(errors.New("DB error"))

	_, err := store.GetPosts()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
