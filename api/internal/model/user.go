package model

import "time"

type User struct {
	Login        string
	Password     string
	Name         string
	Surname      string
	Avatar       string
	IsVarified   bool
	RegisterDate time.Time
	LastSeenDate time.Time
}
