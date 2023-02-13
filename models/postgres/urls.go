package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"urlshortner/models"
)

const URL_KEY_LENGTH = 7

type UrlModel struct {
	DB *sql.DB
}

func (u UrlModel) Insert(url *models.Url) error {
	query := `INSERT INTO urls (target_url, user_id, url_key) VALUES ($1, $2, $3) RETURNING url_key, visits, created_at, updated_at`
	for {
		urlKey := models.GenURLKey(URL_KEY_LENGTH)
		args := []interface{}{url.TargetUrl, url.UserId, urlKey}
		err := u.DB.QueryRow(query, args...).Scan(&url.ShortUrl, &url.Visits, &url.CreatedAt, &url.UpdatedAt)
		if err != nil {
			if err.Error() == `pq: duplicate key value violates unique constraint "urls_pkey"` {
				continue
			} else {
				return err
			}
		}
		return nil
	}
}

func (u UrlModel) GetTargetUrl(urlKey string, userId int64, getInfo bool) (*models.Url, error) {
	var query string
	args := []interface{}{urlKey}
	if !getInfo {
		query = `SELECT * FROM urls WHERE url_key = $1`
		go u.updateVisitsForUrl(urlKey)
	} else {
		args = append(args, userId)
		query = `SELECT * FROM urls WHERE url_key = $1 AND user_id = $2`
	}

	var url models.Url
	err := u.DB.QueryRow(query, args...).Scan(
		&url.ShortUrl,
		&url.TargetUrl,
		&url.Visits,
		&url.UserId,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, models.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &url, nil
}

func (u UrlModel) updateVisitsForUrl(urlKey string) error {
	query := `UPDATE urls SET visits = visits + 1 WHERE url_key = $1`
	result, err := u.DB.Exec(query, urlKey)
	if err != nil {
		fmt.Println("error incrementing visits for url " + urlKey)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}

func (u UrlModel) DeleteUrl(urlKey string, userId int64) error {
	query := `DELETE FROM urls WHERE url_key = $1 AND user_id = $2`
	args := []interface{}{urlKey, userId}
	result, err := u.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}
