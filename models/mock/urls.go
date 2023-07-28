package mock

import (
	"sync"
	"urlshortner/models"
)

const URL_KEY_LENGTH = 7

var rwmutex sync.RWMutex

type UrlModel struct {
	DB map[string]*models.Url
}

func (u UrlModel) Insert(url *models.Url) error {
	/*urlKey := models.GenURLKey(URL_KEY_LENGTH)
	rwmutex.RLock()
	for _, ok := u.DB[urlKey]; ok; {
		urlKey = models.GenURLKey(7)
	}
	rwmutex.RUnlock()

	rwmutex.Lock()
	defer rwmutex.Unlock()
	u.DB[urlKey] = &models.Url{
		TargetUrl: targetUrl,
		ShortUrl:  urlKey,
		UserId:    userId,
		CreatedAt: time.Now(),
		Visits:    0,
	}*/
	return nil
}

func (u UrlModel) ListUrls(userId int64) ([]*models.Url, error) {
	/*urlKey := models.GenURLKey(URL_KEY_LENGTH)
	rwmutex.RLock()
	for _, ok := u.DB[urlKey]; ok; {
		urlKey = models.GenURLKey(7)
	}
	rwmutex.RUnlock()

	rwmutex.Lock()
	defer rwmutex.Unlock()
	u.DB[urlKey] = &models.Url{
		TargetUrl: targetUrl,
		ShortUrl:  urlKey,
		UserId:    userId,
		CreatedAt: time.Now(),
		Visits:    0,
	}*/
	return nil, nil
}

func (u UrlModel) GetTargetUrl(urlKey string, userId int64, getInfo bool) (*models.Url, error) {
	rwmutex.RLock()
	url, ok := u.DB[urlKey]
	rwmutex.RUnlock()
	if !ok {
		return nil, models.ErrRecordNotFound
	}
	if !getInfo {
		go u.updateVisitsForUrl("<url_key>")
	}

	return url, nil
}

func (u UrlModel) updateVisitsForUrl(urlKey string) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	if url, ok := u.DB[urlKey]; ok {
		url.Visits += 1
	}
}

func (u UrlModel) DeleteUrl(urlKey string, userId int64) error {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	delete(u.DB, urlKey)
	return nil
}
