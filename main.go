package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

type NewMessage struct {
	User_id int `json:"user_id"`
	Password string `json:"password"`
	Content string `json:"content"`
}

type Message struct {
	Id int `json:"id"`
	User_id int `json:"user_id"`
	Username string `json:"username"`
	Content string `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	DB_USER := os.Getenv("PYRAMID_CHAT_DB_USER")
	DB_PASS := os.Getenv("PYRAMID_CHAT_DB_PASS")
	DB_HOST := os.Getenv("PYRAMID_CHAT_DB_HOST")
	DB_NAME := os.Getenv("PYRAMID_CHAT_DB_NAME")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/" + DB_NAME + "?parseTime=true"


	r := gin.Default()
	r.GET("/messages", func(c *gin.Context) {
		var (
			message Message
			messages []Message
		)
		// Open DB connection
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer db.Close()

		// make sure our connection is available
		err = db.Ping()
		if err != nil {
			fmt.Print(err.Error())
		}

		rows, err := db.Query("select id, user_id, user_name, content, timestamp FROM messages")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&message.Id, &message.User_id, &message.Username, &message.Content, &message.Timestamp)
			messages = append(messages, message)
			if err != nil {
				log.Fatal(err)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
		})
	})

	r.POST("/messages", func(c *gin.Context) {
		var message NewMessage
		c.BindJSON(&message)

		// Open DB connection
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer db.Close()

		// make sure our connection is available
		err = db.Ping()
		if err != nil {
			fmt.Print(err.Error())
		}

		// Get the username from the DB
		var username string
		var num_chars int
		var password string
		rows, err := db.Query("select username, num_chars, password FROM users where id = ?", message.User_id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		rows.Next()
		err = rows.Scan(&username, &num_chars, &password)
		if err != nil {
			log.Fatal(err)
		}

		if message.Password != password {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		if len(message.Content) > num_chars {
			c.JSON(http.StatusBadRequest, gin.H{"status": "too_many_chars"})
			return
		}

		stmt, err := db.Prepare("insert into messages (user_id, user_name, content) values(?,?,?);")
		_, err = stmt.Exec(message.User_id, username, message.Content)


		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

