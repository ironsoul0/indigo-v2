package db

const (
	USERS_BUCKET = "users"
)

type User struct {
	ChatID    string `json:"chatID"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Activated bool   `json:"activated"`
}
