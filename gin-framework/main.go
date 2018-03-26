package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {

	username := os.Getenv("POSTGRES_USER");
	password := os.Getenv("POSTGRES_PASSWORD")
	hostName := os.Getenv("POSTGRES_HOST")

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		hostName,
		username,
		password,
		"postgres")

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a := App{}
	a.Initialize(db)
	a.Run(":8080")
}
