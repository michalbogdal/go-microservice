package models

import (
	"github.com/jmoiron/sqlx"
)

type Datastore interface {
	Get(id int) (*User, error)
	Update(user *User) error
	Delete(user *User) error
	Create(user *User) error
	List(start, count int) ([]*User, error)
}

type DB struct {
	*sqlx.DB
}

func CreateDB(db *sqlx.DB) *DB {
	return &DB{db}
}
