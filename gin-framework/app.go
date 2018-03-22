package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"go-microservice/gin-framework/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type App struct {
	db     models.Datastore
	router *gin.Engine
}

func (a *App) Initialize(db *sqlx.DB) {
	a.db = models.CreateDB(db)
	a.router = gin.Default()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {

	a.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	a.router.GET("/users", a.getUsers)

	/*a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")*/
}

func (a *App) Run(addr string) {
	a.router.Run(addr)

	log.Fatal(http.ListenAndServe(addr, a.router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	user, err := a.db.Get(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (a *App) getUsers(c *gin.Context) {

	count := c.GetInt("count")
	start := c.GetInt("start")
	/*count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))*/

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := a.db.List(start, count)
	/*if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)*/
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, users)
	}
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := a.db.Create(&user); err != nil {
		fmt.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	user.ID = id
	if err := a.db.Update(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user := models.User{ID: id}
	if err := a.db.Delete(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
