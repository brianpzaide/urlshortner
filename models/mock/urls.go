package mock

import (
	"urlshortner/models"
)

const URL_KEY_LENGTH = 7

type UrlModel struct{}

func (u UrlModel) Insert(url *models.Url) error {

	urlKey := models.GenURLKey(URL_KEY_LENGTH)
	_, ok := DB.Urls[urlKey]
	for ok {
		urlKey = models.GenURLKey(7)
		_, ok = DB.Urls[urlKey]
	}
	DB.Urls[urlKey] = url
	return nil
}

func (u UrlModel) ListUrls(userId int64) ([]*models.Url, error) {
	urlsList := make([]*models.Url, 0)

	for _, url := range DB.Urls {
		if url.UserId == userId {
			temp := url
			urlsList = append(urlsList, temp)
		}
	}
	return urlsList, nil
}

func (u UrlModel) GetTargetUrl(urlKey string, userId int64, getInfo bool) (*models.Url, error) {
	url, ok := DB.Urls[urlKey]
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

func (u UrlModel) DeleteUrl(urlKey string, userId int64) error {

	url, ok := DB.Urls[urlKey]
	if !ok {
		return models.ErrRecordNotFound
	}
	if url.UserId != userId {
		return models.ErrRecordNotFound
	}

	delete(DB.Urls, urlKey)
	return nil
}
