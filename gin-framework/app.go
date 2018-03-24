package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go-microservice/gin-framework/models"
	"log"
	"net/http"
	"strconv"
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

	v1 := a.router.Group("/v1")

	v1.GET("/ping", ping)
	v1.GET("/users", a.getUsers)
	v1.POST("/user", a.createUser)
	v1.GET("/user/:id", a.getUser)
	v1.PUT("/user/:id", a.updateUser)
	v1.DELETE("/user/:id", a.deleteUser)
}

func (a *App) Run(addr string) {
	a.router.Run(addr)

	log.Fatal(http.ListenAndServe(addr, a.router))
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, map[string]string{"error": message})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (a *App) getUser(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	user, err := a.db.Get(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(c, http.StatusNotFound, "User not found")
		default:
			respondWithError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (a *App) getUsers(c *gin.Context) {

	count, err := strconv.Atoi(c.Query("count"))
	start, err2 := strconv.Atoi(c.Query("start"))

	if err != nil && (count > 10 || count < 1) {
		count = 10
	}
	if err2 != nil || start < 0 {
		start = 0
	}

	users, err := a.db.List(start, count)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func (a *App) createUser(c *gin.Context) {
	var user models.User
	decoder := json.NewDecoder(c.Request.Body)

	if err := decoder.Decode(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer c.Request.Body.Close()

	if err := a.db.Create(&user); err != nil {
		fmt.Println(err.Error())
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (a *App) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var user models.User
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer c.Request.Body.Close()
	user.ID = id
	if err := a.db.Update(&user); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *App) deleteUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user := models.User{ID: id}
	if err := a.db.Delete(&user); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{"result": "success"})
}
