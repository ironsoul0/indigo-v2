// Code generated by sqlc. DO NOT EDIT.

package db

import ()

type Grade struct {
	ID         int32  `json:"id"`
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Grade      string `json:"grade"`
	Range      string `json:"range"`
	Percentage string `json:"percentage"`
	CourseName string `json:"course_name"`
}

type User struct {
	ChatID   string `json:"chat_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}
