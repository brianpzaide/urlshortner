package mock

import (
	"urlshortner/models"
)

type UserModel struct {
	DB map[string]*models.User
}

// used for registering new users
func (usr UserModel) Insert(user *models.User) error {

	if usr.DB[user.Email] != nil {
		return models.ErrDuplicateEmail
	}
	usr.DB[user.Email] = user

	return nil
}

// used for users login
func (usr UserModel) GetByEmail(email string) (*models.User, error) {

	user, ok := usr.DB[email]
	if !ok {
		return nil, models.ErrRecordNotFound
	}
	return user, nil
}
