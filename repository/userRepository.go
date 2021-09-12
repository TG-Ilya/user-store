package repository

import (
	"fmt"
	"user-store/internal/models"
	"user-store/storage/sqlite"
)

type UserHandler interface {
	CreateUser(db *sqlite.Storage, user *models.User) (*models.User, error)
	GetAllUser(db *sqlite.Storage) ([]*models.User, error)
	UpdateUser(db *sqlite.Storage, id int, user *models.User) (*models.User, error)
	DeleteUser(db *sqlite.Storage, id int) error
	GetUser(db *sqlite.Storage, id int) (*models.User, error)
}

type userRepository struct{}

func CreateUserHandler() *userRepository {
	return &userRepository{}
}

//Create User in db
func (ur *userRepository) CreateUser(db *sqlite.Storage, usr *models.User) (*models.User, error) {

	query := fmt.Sprintf("INSERT INTO users (name, birth_date) VALUES ($1, $2) RETURNING id")
	err := db.Sql.QueryRow(query, usr.Name, usr.BirthDate).Scan(&usr.ID)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

//Get all users from db
func (ur *userRepository) GetAllUsers(db *sqlite.Storage) ([]*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM users")
	rows, err := db.Sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		usr := models.User{}
		_ = rows.Scan(
			&usr.ID,
			&usr.Name,
			&usr.BirthDate,
		)
		users = append(users, &usr)
	}
	return users, nil
}

//Get user by id
func (ur *userRepository) GetUser(db *sqlite.Storage, id int64) (*models.User, error) {
	usr := &models.User{}
	query := fmt.Sprintf("SELECT * FROM users WHERE id = $1")

	if err := db.Sql.QueryRow(query, id).Scan(&usr.ID, &usr.Name, &usr.BirthDate); err != nil {
		return nil, err
	}
	return usr, nil
}

//Delete user from db
func (ur *userRepository) DeleteUser(db *sqlite.Storage, id int64) error {
	query := fmt.Sprintf("DELETE FROM users WHERE id = $1")
	_, err := db.Sql.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

//Update user from db
func (ur *userRepository) UpdateUser(db *sqlite.Storage, id int64, usr *models.User) (*models.User, error) {
	query := fmt.Sprintf("UPDATE users SET name = $1, birth_date = $2 where id = $3 RETURNING id, name, birth_date")
	if err := db.Sql.QueryRow(query, usr.Name, usr.BirthDate, id).Scan(&usr.ID, &usr.Name, &usr.BirthDate); err != nil {
		return nil, err
	}
	return usr, nil
}
