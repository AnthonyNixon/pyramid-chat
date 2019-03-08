package users

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/AnthonyNixon/pyramid-chat/values"

	"github.com/AnthonyNixon/pyramid-chat/database"
	"github.com/AnthonyNixon/pyramid-chat/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// https://www.sohamkamani.com/blog/2018/02/25/golang-password-authentication-and-storage/
	var newUser types.NewUser
	c.BindJSON(&newUser)

	unique, err := IsNewUserUnique(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email is already taken"})
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	referral_id, err := LookUpUserID(newUser.Referral_code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	stmt, err := db.Prepare("insert into users (username, email, password, invited_by, num_chars, referral_code) values(?,?,?,?,?,?);")
	_, err = stmt.Exec(newUser.Username, newUser.Email, hashedPassword, referral_id, values.STARTING_CHARS, NewReferralCode())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if referral_id != 0 {
		IncreaseNumChars(referral_id)
	}

}

func LookUpUserID(referral_code string) (int, error) {
	if referral_code == "" {
		return 0, nil
	}
	db, err := database.GetConnection()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	rows, err := db.Query("select id FROM users where referral_code = ?", referral_code)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var user_id int
	rows.Next()
	err = rows.Scan(&user_id)
	if err != nil {
		return -1, err
	}

	return user_id, nil
}

func IsNewUserUnique(newUser types.NewUser) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("select COUNT(*) FROM users where username = ? OR email = ?", newUser.Username, newUser.Email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func IncreaseNumChars(user_id int) error {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("update users set num_chars = num_chars + ? where id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(values.REFERRAL_CHAR_INCREASE_FIRST_LEVEL, user_id)
	if err != nil {
		return err
	}

	return nil
}

func NewReferralCode() string {
	var runes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, values.DEFAULT_REFERRAL_CODE_LEN)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}
