package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	c.IndentedJSON(http.StatusOK, users)
}

func getUserById(id int) (*user, error) {
	for i, u := range users {
		if u.User_id == id {
			return &users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func getUser(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := getUserById(userId)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, user)
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
