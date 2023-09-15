package mock

import (
	"log"
	"urlshortner/models"
)

const URL_KEY_LENGTH = 7

type UrlModel struct {
	db *MockDB
}

func NewUrlModel(db *MockDB) *UrlModel {
	return &UrlModel{db: db}
}

func (u *UrlModel) Insert(url *models.Url) error {

	urlKey := models.GenURLKey(URL_KEY_LENGTH)
	_, ok := u.db.Urls[urlKey]
	for ok {
		urlKey = models.GenURLKey(7)
		_, ok = u.db.Urls[urlKey]
	}
	u.db.Urls[urlKey] = url
	return nil
}

func (u *UrlModel) ListUrls(userId int64) ([]*models.Url, error) {
	urlsList := make([]*models.Url, 0)

	for _, url := range u.db.Urls {
		if url.UserId == userId {
			temp := url
			urlsList = append(urlsList, temp)
		}
	}
	return urlsList, nil
}

func (u *UrlModel) GetTargetUrl(urlKey string, userId int64, getInfo bool) (*models.Url, error) {
	url, ok := u.db.Urls[urlKey]
	if !ok {
		return nil, models.ErrRecordNotFound
	}
	if !getInfo {
		url.Visits = url.Visits + 1
		return url, nil
	} else if userId == url.UserId {
		return url, nil
	} else {
		return nil, models.ErrRecordNotFound
	}

}

func (u *UrlModel) DeleteUrl(urlKey string, userId int64) error {
	log.Printf("DeleteUrl userid: %d, urlKey: %s", userId, urlKey)
	url, ok := u.db.Urls[urlKey]
	if !ok {
		return models.ErrRecordNotFound
	}
	if url.UserId != userId {
		return models.ErrRecordNotFound
	}

	delete(u.db.Urls, urlKey)
	return nil
}
