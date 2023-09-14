package postgres

import (
	"testing"
	"time"
	"urlshortner/models"
)

func Test_repo(t *testing.T) {
	testUser := &models.User{
		Email: "abc@example.com",
	}
	testUser.Password.Set("12345")

	// testing for UserModel.Insert
	err := testUserRepo.Insert(testUser)
	if err != nil {
		t.Errorf("Usermodel.Insert returned an error: %s", err)
	}

	// testing for UserModel.GetByEmail
	fetchedUser, err := testUserRepo.GetByEmail("abc@example.com")
	if err != nil {
		t.Errorf("UserModel.GetByEmail returned an error: %s", err)
	}

	if fetchedUser.ID < 1 {
		t.Errorf("UserModel.GetByEmail returned wrong userId; expected >= 1, but got %d", fetchedUser.ID)
	}

	// testing for TokenModel.New, TokenModel.Insert
	token, err := testTokenRepo.New(fetchedUser.ID, 1*time.Hour, models.ScopeAuthentication)
	if err != nil {
		t.Error("TokenModel.New, TokenModel.Insert returned error", err)
	}

	// testing for UserModel.GetForToken
	fetchedUserByToken, err := testUserRepo.GetForToken(token.Scope, token.Plaintext)
	if err != nil {
		t.Error("UserModel.GetForToken returned error", err)
	}
	if fetchedUserByToken.ID != fetchedUser.ID || fetchedUserByToken.Email != fetchedUser.Email {
		t.Errorf("UserModel.GetForToken fetched wrong user %d-%s, expected %d-%s", fetchedUser.ID, fetchedUser.Email, fetchedUserByToken.ID, fetchedUser.Email)
	}

	// testing for UrlModel.Insert
	urlInfo := models.Url{
		TargetUrl: "https://stackoverflow.com/",
		UserId:    fetchedUser.ID,
	}
	err = testUrlRepo.Insert(&urlInfo)
	if err != nil {
		t.Error("UrlModel.Insert returned error", err)
	}

	// testing for UrlModel.ListUrls
	urls, err := testUrlRepo.ListUrls(fetchedUser.ID)
	if err != nil {
		t.Error("UrlModel.ListUrls returned error", err)
	}
	if len(urls) != 1 {
		t.Errorf("UrlModel.ListUrls returned error user created only 1 url but method returned %d", len(urls))
	}
	if urls[0].TargetUrl != urlInfo.TargetUrl {
		t.Errorf("UrlModel.ListUrls user created only 1 url for %s but method returned url for %s", urlInfo.TargetUrl, urls[0].TargetUrl)
	}

	// testing for UrlModel.GetTargetUrl, UrlModel.updateVisitsForUrl
	targetUrl := urls[0].TargetUrl
	short_url := urls[0].ShortUrl
	userId := urls[0].UserId
	oldVisits := urls[0].Visits
	// testing for UrlModel.GetTargetUrl with getInfo=false and correct url_key
	fetchedUrl, err := testUrlRepo.GetTargetUrl(short_url, userId, false)
	if err != nil {
		t.Fatal("UrlModel.GetTargetUrl returned error", err)
	}
	if fetchedUrl.TargetUrl != targetUrl || fetchedUrl.UserId != userId {
		t.Errorf("UrlModel.GetTargetUrl returned wrong url, method returned %d-%s, expected %d-%s", fetchedUrl.UserId, fetchedUrl.TargetUrl, userId, targetUrl)
	}

	// testing for UrlModel.GetTargetUrl with getInfo=false and incorrect url_key
	fetchedUrl, err = testUrlRepo.GetTargetUrl("nonexisting url_key", userId, false)
	if err == nil {
		t.Error("UrlModel.GetTargetUrl returned url for nonexisting url_key")
	}

	// testing for UrlModel.GetTargetUrl with getInfo=true and incorrect userId
	fetchedUrl, err = testUrlRepo.GetTargetUrl(short_url, 0, true)
	if err == nil {
		t.Error("UrlModel.GetTargetUrl returned url for nonexisting userId=0")
	}

	time.Sleep(1 * time.Second)

	// testing for UrlModel.GetTargetUrl with getInfo=true and correct userId and correct url_key
	fetchedUrl1, err := testUrlRepo.GetTargetUrl(short_url, userId, true)
	if err != nil {
		t.Fatal("UrlModel.GetTargetUrl returned error", err)
	}
	if fetchedUrl1.Visits != oldVisits+1 {
		t.Errorf("UrlModel.updateVisitsForUrl error visits did not get updated old_visits=%d new_visits=%d", oldVisits, fetchedUrl1.Visits)
	}

	// testing for UrlModel.GetTargetUrl with getInfo=true and correct userId and incorrect url_key
	fetchedUrl, err = testUrlRepo.GetTargetUrl("nonexisting url_key", userId, true)
	if err == nil {
		t.Error("UrlModel.GetTargetUrl fetched url for non existing key")
	}

	// testing for UrlModel.DeleteUrl with correct userId and incorrect url_key
	err = testUrlRepo.DeleteUrl("nonexisting url_key", userId)
	if err == nil {
		t.Error("UrlModel.DeleteUrl deleted url for non existing key")
	}

	// testing for UrlModel.DeleteUrl with incorrect userId and correct url_key
	err = testUrlRepo.DeleteUrl(short_url, 0)
	if err == nil {
		t.Error("UrlModel.DeleteUrl deleted url for non existing userId=0 ")
	}

	// testing for UrlModel.DeleteUrl with correct userId and correct url_key
	err = testUrlRepo.DeleteUrl(short_url, userId)
	if err != nil {
		t.Error("UrlModel.DeleteUrl returned error", err)
	}

	// testing for TokenModel.DeleteAllForUser
	err = testTokenRepo.DeleteAllForUser(models.ScopeAuthentication, fetchedUser.ID)
	if err != nil {
		t.Fatal("TokenModel.DeleteAllForUser, returned error", err)
	}

}
