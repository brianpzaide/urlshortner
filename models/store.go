package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const ScopeAuthentication = "authentication"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

var AnonymousUser = &User{}

type password struct {
	Plaintext *string
	Hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type Url struct {
	TargetUrl string    `json:"target_url"`
	ShortUrl  string    `json:"url_key"`
	Visits    int64     `json:"visits"`
	UserId    int64     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type UserModelInterface interface {
	Insert(*User) error
	GetByEmail(string) (*User, error)
	GetForToken(string, string) (*User, error)
}
type URLModelInterface interface {
	Insert(*Url) error
	ListUrls(int64) ([]*Url, error)
	GetTargetUrl(string, int64, bool) (*Url, error)
	DeleteUrl(string, int64) error
}
type TokenModelInterface interface {
	New(int64, time.Duration, string) (*Token, error)
	Insert(*Token) error
	DeleteAllForUser(string, int64) error
}

type Models struct {
	Users  UserModelInterface
	Urls   URLModelInterface
	Tokens TokenModelInterface
}

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenURLKey(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}

func GenerateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	fmt.Println("no errors till store 117")
	randomBytes := make([]byte, 16)
	fmt.Println("no errors till store 119")
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("no errors till store 123")
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}
