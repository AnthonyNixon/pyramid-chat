package main

import (
	"github.com/AnthonyNixon/pyramid-chat/database"
	"github.com/AnthonyNixon/pyramid-chat/messages"
	"github.com/AnthonyNixon/pyramid-chat/users"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.Initialize()
}

func main() {
	r := gin.Default()
	r.GET("/messages", func(c *gin.Context) {
		messages.GetAll(c)
	})

	r.POST("/messages", func(c *gin.Context) {
		messages.New(c)
	})

	r.POST("/signup", func(c *gin.Context) {
		users.SignUp(c)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
