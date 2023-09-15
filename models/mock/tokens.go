package mock

import (
	"time"
	"urlshortner/models"
)

type TokenModel struct {
	db *MockDB
}

func NewTokenModel(db *MockDB) *TokenModel {
	return &TokenModel{db: db}
}

func (t *TokenModel) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := models.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

func (t *TokenModel) Insert(token *models.Token) error {
	t.db.Tokens[token.Plaintext] = token
	return nil
}

func (t *TokenModel) DeleteAllForUser(scope string, userId int64) error {
	toBeDeleted := make([]string, 0)
	for key, value := range t.db.Tokens {
		if value.UserId == userId && value.Scope == scope {
			toBeDeleted = append(toBeDeleted, key)
		}
	}
	for _, key := range toBeDeleted {
		delete(t.db.Tokens, key)
	}
	return nil
}
