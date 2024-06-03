package postgres_test

import (
	"database/sql"
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

		_, err := store.AddPost("Test Title", "Test Content", true)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestAddComment(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	t.Run("AddCommentToPost", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "allow_comments"}).AddRow("post-id", "Test Title", "Test Content", time.Now(), true)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at, allow_comments FROM posts WHERE id = $1")).
			WithArgs("post-id").WillReturnRows(row)

		row = sqlmock.NewRows([]string{"id"}).AddRow("post-id")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM comments WHERE post_id = $1")).
			WithArgs("post-id").WillReturnRows(row)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO comments (id, post_id, parent_comment_id, content, created_at) VALUES ($1, $2, NULL, $3, $4)")).
			WithArgs(sqlmock.AnyArg(), "post-id", "Test Comment", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		_, err := store.AddComment("post-id", "", "Test Comment")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetPosts(t *testing.T) {
	db, mock, _ := NewMock()
	defer db.Close()

	store := postgres.DataStorePostgres{DB: db}
	storage.DataBase = &store

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "true"}).
		AddRow("post-id", "Test Title", "Test Content", time.Now(), "true")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, content, created_at, allow_comments FROM posts")).WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM comments WHERE post_id = $1")).WithArgs("post-id").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("comment-id"))

	_, err := store.GetPosts()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
