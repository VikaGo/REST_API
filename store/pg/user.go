package pg

import (
	"context"
	"database/sql"
	"github.com/VikaGo/REST_API/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserRepo ...
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetUser retrieves user from Postgres
func (repo *UserRepo) GetUser(ctx context.Context, id uuid.UUID) (*model.DBUser, error) {
	user := &model.DBUser{}
	err := repo.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// CreateUser creates user in Postgres
func (repo *UserRepo) CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	_, err := repo.db.NamedExec("INSERT INTO users (id, firstname, lastname, nickname, password) VALUES (:id, :firstname, :lastname, :nickname, :password)", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates user in Postgres
func (repo *UserRepo) UpdateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	_, err := repo.db.NamedExec("UPDATE users SET id = :id, firstname = :firstname, lastname =:lastname, nickname =:nickname, password =:password WHERE id = :id", user)
	if err != nil {
		if err == sql.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes user in Postgres
func (repo *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (repo *UserRepo) GetPassword(ctx context.Context, id uuid.UUID) (string, error) {
	var password string
	err := repo.db.Get(&password, "SELECT password FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows { // If the user is not found, return an empty string and no error.
			return "", nil
		}
		return "", err // Otherwise, return an error if one occurred during the database query.
	}

	// If the query was successful, return the retrieved password as a string and no error.
	return password, nil
}

func (repo *UserRepo) GetUserByNickname(ctx context.Context, nickname string) (*model.DBUser, error) {
	user := &model.DBUser{}
	err := repo.db.Get(user, "SELECT * FROM users WHERE nickname = $1", nickname)
	if err != nil {
		if err == sql.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
