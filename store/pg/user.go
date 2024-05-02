package pg

import (
	"context"
	"github.com/VikaGo/REST_API/model"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// UserPgRepo ...
type UserPgRepo struct {
	db *DB
}

// NewUserRepo ...
func NewUserRepo(db *DB) *UserPgRepo {
	return &UserPgRepo{db: db}
}

// GetUser retrieves user from Postgres
func (repo *UserPgRepo) GetUser(ctx context.Context, id uuid.UUID) (*model.DBUser, error) {
	user := &model.DBUser{}
	err := repo.db.Model(user).
		Where("id = ?", id).
		Select()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// CreateUser creates user in Postgres
func (repo *UserPgRepo) CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	_, err := repo.db.Model(user).
		Returning("*").
		Insert()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates user in Postgres
func (repo *UserPgRepo) UpdateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	_, err := repo.db.Model(user).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes user in Postgres
func (repo *UserPgRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Model((*model.DBUser)(nil)).
		Where("id = ?", id).
		Delete()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (repo *UserPgRepo) GetPassword(ctx context.Context, id uuid.UUID) (string, error) {
	// Create a user model object to store the result.
	user := &model.DBUser{}

	// Execute a database query to select the user's password with the specified ID.
	err := repo.db.Model(user).
		Column("password").
		Where("id = ?", id).
		Select()

	if err != nil {
		if err == pg.ErrNoRows { // If the user is not found, return an empty string and no error.
			return "", nil
		}
		return "", err // Otherwise, return an error if one occurred during the database query.
	}

	// If the query was successful, return the retrieved password as a string and no error.
	return user.Password, nil
}
