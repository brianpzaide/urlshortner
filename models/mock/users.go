package mock

import (
	"time"
	"urlshortner/models"
)

type MockDB struct {
	Users       map[int64]*models.User
	EmailLookup map[string]int64
	Urls        map[string]*models.Url
	Tokens      map[string]*models.Token
}

var (
	Count int64
	DB    MockDB
)

type UserModel struct{}

// used for registering new users
func (usr UserModel) Insert(user *models.User) error {

	if _, ok := DB.EmailLookup[user.Email]; ok {
		return models.ErrDuplicateEmail
	}
	DB.Users[int64(Count)] = user
	DB.Users[int64(Count)].ID = int64(Count)
	DB.EmailLookup[user.Email] = int64(Count)
	Count = Count + 1

	return nil
}

// used for users login
func (usr UserModel) GetByEmail(email string) (*models.User, error) {

	userId, ok := DB.EmailLookup[email]
	if !ok {
		return nil, models.ErrRecordNotFound
	}
	return DB.Users[userId], nil
}

func (u UserModel) GetForToken(tokenScope, tokenPlaintext string) (*models.User, error) {
	token, ok := DB.Tokens[tokenPlaintext]
	if !ok {
		return nil, models.ErrRecordNotFound
	}

	if time.Now().After(token.Expiry) {
		return nil, models.ErrRecordNotFound
	}
	return DB.Users[token.UserId], nil
}
