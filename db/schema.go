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

type SceneID int

const (
	NO_SCENE           SceneID = iota
	REGISTER_SCENE             = iota
	VERIFICATION_SCENE         = iota
)

type RegisterMeta struct {
	Username string
	Password string
}

type VerificationMeta struct {
	Code string
}

type Scene struct {
	ChatID int
	Scene  SceneID
	RegisterMeta
	VerificationMeta
	Step int
}
