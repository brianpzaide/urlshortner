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

type UserModel struct {
	db    *MockDB
	count int
}

func NewUserModel(db *MockDB) *UserModel {
	return &UserModel{db: db, count: 0}
}

// used for registering new users
func (usr *UserModel) Insert(user *models.User) error {

	if _, ok := usr.db.EmailLookup[user.Email]; ok {
		return models.ErrDuplicateEmail
	}
	usr.db.Users[int64(usr.count)] = user
	usr.db.Users[int64(usr.count)].ID = int64(usr.count)
	usr.db.EmailLookup[user.Email] = int64(usr.count)
	usr.count = usr.count + 1

	return nil
}

// used for users login
func (usr *UserModel) GetByEmail(email string) (*models.User, error) {

	userId, ok := usr.db.EmailLookup[email]
	if !ok {
		return nil, models.ErrRecordNotFound
	}
	return usr.db.Users[userId], nil
}

func (usr *UserModel) GetForToken(tokenScope, tokenPlaintext string) (*models.User, error) {
	token, ok := usr.db.Tokens[tokenPlaintext]
	if !ok {
		return nil, models.ErrRecordNotFound
	}

	if time.Now().After(token.Expiry) {
		return nil, models.ErrRecordNotFound
	}
	return usr.db.Users[token.UserId], nil
}
