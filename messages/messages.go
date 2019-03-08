package messages

import (
	"net/http"

	"github.com/AnthonyNixon/pyramid-chat/auth"

	"github.com/AnthonyNixon/pyramid-chat/database"
	"github.com/AnthonyNixon/pyramid-chat/types"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	var (
		message  types.Message
		messages []types.Message
	)
	// Open DB connection
	database, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer database.Close()

	rows, err := database.Query("select id, user_id, user_name, content, timestamp FROM messages")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&message.Id, &message.User_id, &message.Username, &message.Content, &message.Timestamp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		messages = append(messages, message)
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func New(c *gin.Context) {
	var message types.NewMessage
	c.BindJSON(&message)

	// Open DB connection
	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	// Get the username from the DB
	var username string
	var num_chars int
	rows, err := db.Query("select username, num_chars FROM users where id = ?", message.User_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&username, &num_chars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	authed, err := auth.IsAuthed(username, message.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !authed {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	if len(message.Content) > num_chars {
		c.JSON(http.StatusBadRequest, gin.H{"status": "too_many_chars"})
		return
	}

	stmt, err := db.Prepare("insert into messages (user_id, user_name, content) values(?,?,?);")
	_, err = stmt.Exec(message.User_id, username, message.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}
