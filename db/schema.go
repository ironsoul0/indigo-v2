package db

import "github.com/ironsoul0/indigo-v2/scrapers/moodle"

const (
	USERS_BUCKET  = "users"
	STATES_BUCKET = "states"
)

type User struct {
	ChatID    int             `json:"chatID"`
	Username  string          `json:"username"`
	Password  string          `json:"password"`
	Activated bool            `json:"activated"`
	Courses   []moodle.Course `json:"courses"`
}
