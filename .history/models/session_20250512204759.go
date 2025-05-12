package models

import "time"

type Session struct {
	UserID      int       `json:"user_id"`
	AccessToken string    `json:"access_token"`
	LoginTime   time.Time `json:"login_time"`
}

type RevokedToken struct {
	JTI string `json:"jti"`
}

type Log struct {
	FamilyID          int       `json:"familyID"`
	UserID            int       `json:"userID"`
	ChangeDescription string    `json:"changeDescription"`
	Timestamp         time.Time `json:"timestamp"`
}
