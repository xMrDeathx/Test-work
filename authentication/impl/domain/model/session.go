package model

import "github.com/google/uuid"

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
	UserIP string
	Token  RefreshToken
}
