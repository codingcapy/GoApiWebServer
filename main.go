package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgresql://postgres:Stevefox1!@localhost/capychatdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
}

type user struct {
	User_id      int    `json:"user_id`
	Username     string `json:"username`
	Password     string `json:"password`
	Email        string `json:"email`
	Display_name string `json:"display_name`
	Created_at   string `json:"created_at`
	Active       bool   `json:"active`
}

var users = []user{
	{User_id: 1, Username: "asdf", Password: "asdf", Email: "asdf@a.a", Display_name: "asdf", Created_at: "April 15, 2024 3:51PM", Active: true},
}

func getHome(c *gin.Context) {
	c.JSON(http.StatusOK, "welcome")
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT user_id, username, password, email, display_name, created_at, active FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []*user
	for rows.Next() {
		var u user
		err := rows.Scan(&u.User_id, &u.Username, &u.Password, &u.Email, &u.Display_name, &u.Created_at, &u.Active)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	rows, err := db.Query("SELECT * FROM users WHERE user_id = $1", userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var users []*user
	for rows.Next() {
		var u user
		err := rows.Scan(&u.User_id, &u.Username, &u.Password, &u.Email, &u.Display_name, &u.Created_at, &u.Active)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func main() {
	router := gin.Default()
	router.GET("/", getHome)
	router.GET("/users", getUsers)
	router.GET("/users/:userId", getUser)
	router.POST("/users", createUser)
	router.Run("localhost:3434")
}
