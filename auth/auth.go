package auth

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/AnthonyNixon/pyramid-chat/database"
)

func IsAuthed(username string, password string) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var storedPassword string
	result := db.QueryRow("select password FROM users where username = ?", username)
	if err != nil {
		return false, err
	}

	err = result.Scan(&storedPassword)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, nil
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}
