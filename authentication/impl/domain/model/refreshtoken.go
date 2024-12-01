package model

type RefreshToken struct {
	Token     []byte
	ExpiresIn int64
}
