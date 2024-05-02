package model

import (
	"time"

	"github.com/google/uuid"
)

// User is a JSON user
type User struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role" validate:"required"`
	Firstname string    `json:"firstname" validate:"required"`
	Lastname  string    `json:"lastname" validate:"required"`
	Nickname  string    `json:"nickname" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// ToDB converts User to DBUser
func (user *User) ToDB() *DBUser {

	if user == nil {
		return nil
	}

	return &DBUser{
		ID:        user.ID,
		Role:      user.Role,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Nickname:  user.Nickname,
		Password:  user.Password,
	}
}

// DBUser is a Postgres user
type DBUser struct {
	ID        uuid.UUID `db:"id"`
	Role      string    `db:"role"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Nickname  string    `db:"nickname"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

// ToWeb converts DBUser to User
func (dbUser *DBUser) ToWeb() *User {

	if dbUser == nil {
		return nil
	}

	return &User{
		ID:        dbUser.ID,
		Role:      dbUser.Role,
		Firstname: dbUser.Firstname,
		Lastname:  dbUser.Lastname,
		Nickname:  dbUser.Nickname,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		DeletedAt: dbUser.DeletedAt,
	}
}
