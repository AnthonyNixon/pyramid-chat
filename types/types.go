package types

import "time"

type NewMessage struct {
	User_id  int    `json:"user_id"`
	Password string `json:"password"`
	Content  string `json:"content"`
}

type Message struct {
	Id        int       `json:"id"`
	User_id   int       `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type NewUser struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Referral_code string `json:"referral_code"`
}

type UserAuth struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
