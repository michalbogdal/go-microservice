package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func main() {

	configFile := flag.String("configFile", "../resources/config.json", "config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.ReadInConfig()

	username := viper.GetString("database.user")
	password := viper.GetString("database.password")
	hostName := viper.GetString("database.hostname")
	appPort := viper.GetString("host.port")

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		hostName,
		username,
		password,
		"postgres")

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("starting app on port", appPort)
	a := App{}
	a.Initialize(db)
	a.Run(":" + appPort)
}
