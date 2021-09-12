package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"log"
	"regexp"
	"testing"
	"user-store/internal/models"
	"user-store/storage/sqlite"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error '%s' was not expected when opening database connection", err)
	}

	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	storage := sqlite.Storage{Sql: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO users \\(name, birth_date\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WillReturnRows(rows)

	ur := CreateUserHandler()
	us, err := ur.CreateUser(&storage, &models.User{Name: swag.String("Ivan"), BirthDate: swag.String("12-05-1980")})
	if err != nil {
		t.Errorf("Error '%s' was not expected, while inserting a row", err)
	}

	assert.NotEmpty(t, us)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	storage := sqlite.Storage{Sql: db}

	rows := sqlmock.NewRows([]string{"id", "mane", "birth_date"}).AddRow(1, "Ivan", "12-05-1980")
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").WillReturnRows(rows)

	ur := CreateUserHandler()
	us, err := ur.GetUser(&storage, 1)
	if err != nil {
		t.Errorf("Error '%s' was not expected, while selecting a row", err)
	}

	assert.NotEmpty(t, us)
	assert.NoError(t, err)
}

func TestGetAllUsers(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	storage := sqlite.Storage{Sql: db}

	rows := sqlmock.NewRows([]string{"id", "mane", "birth_date"}).
		AddRow(1, "Ivan", "12-05-1980").
		AddRow(2, "Petr", "12-05-1981")
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	ur := CreateUserHandler()
	us, err := ur.GetAllUsers(&storage)
	if err != nil {
		t.Errorf("Error '%s' was not expected, while selecting rows", err)
	}

	assert.NotEmpty(t, us)
	assert.NoError(t, err)
	assert.Len(t, us, 2)
}

func TestUpdateUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	storage := sqlite.Storage{Sql: db}

	rows := sqlmock.NewRows([]string{"id", "mane", "birth_date"}).AddRow(1, "Petr", "12-05-1980")
	mock.ExpectQuery("UPDATE users SET name = \\$1, birth_date = \\$2 where id = \\$3 RETURNING id, name, birth_date").
		WillReturnRows(rows)

	ur := CreateUserHandler()
	us, err := ur.UpdateUser(&storage, 1, &models.User{Name: swag.String("Petr"), BirthDate: swag.String("12-05-1980")})
	if err != nil {
		t.Errorf("Error '%s' was not expected, while updating row", err)
	}

	assert.NotEmpty(t, us)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	storage := sqlite.Storage{Sql: db}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	ur := CreateUserHandler()
	err := ur.DeleteUser(&storage, 1)
	if err != nil {
		t.Errorf("Error '%s' was not expected, while updating row", err)
	}

	assert.NoError(t, err)
}
