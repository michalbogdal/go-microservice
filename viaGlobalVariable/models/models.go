package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

var db *sqlx.DB

func InitializeDb(database *sqlx.DB) {
	db = database
}

func (u *User) Get() error {
	return db.Get(u, "SELECT name, email from users Where id=$1", u.ID)
}

func (u *User) Update() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4", u.Name, u.Email, string(hashedPassword), u.ID)
	return err
}

func (u *User) Delete() error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)
	return err
}

func (u *User) Create() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.QueryRow("INSERT INTO users(name, email, password) VALUES($1,$2,$3) RETURNING id", u.Name, u.Email, string(hashedPassword)).Scan(&u.ID)
}

func List(start, count int) ([]User, error) {
	users := []User{}
	err := db.Select(&users, "SELECT id, name, email FROM users LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}
	return users, nil
}
