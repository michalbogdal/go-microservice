package models

import (
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func (db *DB) Get(id int) (*User, error) {
	var u User
	err := db.DB.Get(&u, "SELECT id, name, email from users Where id=$1", id)
	return &u, err
}

func (db *DB) Update(u *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	_, err = db.DB.Exec("UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4", u.Name, u.Email, string(hashedPassword), &u.ID)
	return err
}

func (db *DB) Delete(u *User) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)
	return err
}

func (db *DB) Create(u *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.QueryRow("INSERT INTO users(name, email, password) VALUES($1,$2,$3) RETURNING id", &u.Name, &u.Email, string(hashedPassword)).Scan(&u.ID)
}

func (db *DB) List(start, count int) ([]*User, error) {
	users := []*User{}
	err := db.Select(&users, "SELECT id, name, email FROM users LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}
	return users, nil
}
