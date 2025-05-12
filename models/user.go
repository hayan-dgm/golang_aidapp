package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Hidden from JSON output
	IsAdmin  bool   `json:"isAdmin"`
}
