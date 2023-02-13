package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"
	"urlshortner/models"
)

type UserModel struct {
	DB *sql.DB
}

// used for registering new users
func (u UserModel) Insert(user *models.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	args := []interface{}{user.Email, user.Password.Hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return models.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// used for users login
func (u UserModel) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`

	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, models.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) GetForToken(tokenScope, tokenPlaintext string) (*models.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `SELECT id, email, password_hash, created_at
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.token_hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, models.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
