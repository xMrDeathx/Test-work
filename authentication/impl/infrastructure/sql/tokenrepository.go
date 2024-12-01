package sql

import (
	"TestWork/authentication/impl/domain/model"
	"TestWork/authentication/impl/domain/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewTokenStorage(conn *pgxpool.Pool) repositories.TokenRepository {
	return &tokenRepository{conn: conn}
}

type tokenRepository struct {
	conn *pgxpool.Pool
}

// Получение refresh token, записанного в базу данных
func (repo *tokenRepository) GetToken(context context.Context, userID uuid.UUID) (model.RefreshToken, string, error) {
	var token model.RefreshToken
	var userIP string

	err := repo.conn.QueryRow(context, `
		SELECT token, expires_in, user_ip FROM user_token
		WHERE user_id = $1
	`, userID).Scan(&token.Token, &token.ExpiresIn, &userIP)

	if err == pgx.ErrNoRows {
		return model.RefreshToken{}, "", errors.New("user with token not found")
	} else if err != nil {
		return model.RefreshToken{}, "", err
	}

	return token, userIP, nil
}

// Обновление refresh token в базе данных
func (repo *tokenRepository) UpdateToken(context context.Context, oldToken []byte, newToken model.RefreshToken, requestIP string) error {
	_, err := repo.conn.Exec(context, `
		UPDATE user_token
		SET token = $1, user_ip = $2, expires_in = $3
		WHERE token = $4
	`, newToken.Token, requestIP, newToken.ExpiresIn, oldToken)

	return err
}

// Сохранение refresh token в базе данных
func (repo *tokenRepository) SaveToken(context context.Context, session model.Session) error {
	_, err := repo.conn.Exec(context, `
		INSERT INTO user_token (id, user_id, user_ip, token, expires_in)
		VALUES ($1, $2, $3, $4, $5)
	`, session.ID, session.UserID, session.UserIP, session.Token.Token, session.Token.ExpiresIn)

	return err
}
