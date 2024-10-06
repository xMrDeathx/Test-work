package sql

import (
	"TestWork/authentication/impl/domain/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func NewTokenStorage(conn *pgxpool.Pool) repositories.TokenRepository {
	return &tokenRepository{conn: conn}
}

type tokenRepository struct {
	conn *pgxpool.Pool
}

// Проверка валидности refresh token полученного в запросе
func (repo *tokenRepository) ValidateRefreshToken(context context.Context, userID uuid.UUID, oldToken uuid.UUID) (bool, error) {
	var hashedOldToken []byte

	err := repo.conn.QueryRow(context, `
		SELECT token FROM user_token
		WHERE userId = $1
	`, userID).Scan(&hashedOldToken)

	if err == pgx.ErrNoRows {
		return false, errors.New("user with token not found")
	} else if err != nil {
		return false, err
	}

	err = validateToken(oldToken, hashedOldToken)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Обновление refresh token в базе данных
func (repo *tokenRepository) UpdateToken(context context.Context, token uuid.UUID, userID uuid.UUID) error {
	hashedToken, err := hashRefreshToken(token)
	if err != nil {
		return err
	}

	_, err = repo.conn.Exec(context, `
		UPDATE user_token
		SET token = $1
		WHERE userid = $2
	`, hashedToken, userID)

	return err
}

// Сохранение refresh token в базе данных
func (repo *tokenRepository) SaveToken(context context.Context, token uuid.UUID, userID uuid.UUID) error {
	hashedToken, err := hashRefreshToken(token)
	if err != nil {
		return err
	}

	_, err = repo.conn.Exec(context, `
		INSERT INTO user_token (id, userid, token)
		VALUES ($1, $2, $3)
	`, uuid.New(), userID, hashedToken)

	return err
}

func hashRefreshToken(token uuid.UUID) ([]byte, error) {
	byteToken, err := token.MarshalBinary()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(byteToken, 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func validateToken(oldToken uuid.UUID, hashedOldToken []byte) error {
	byteOldToken, err := oldToken.MarshalBinary()
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(hashedOldToken, byteOldToken)
}
