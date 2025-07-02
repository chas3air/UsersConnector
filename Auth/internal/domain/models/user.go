package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}
