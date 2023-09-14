package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const dsn = "host=localhost port=5432 user=postgres password=mysecretpassword dbname=urlshortner sslmode=disable timezone=UTC connect_timeout=5"

var (
	testDB        *sql.DB
	testUserRepo  UserModel
	testUrlRepo   UrlModel
	testTokenRepo TokenModel
)

func TestMain(m *testing.M) {

	testDB, err := openDB()
	if err != nil {
		log.Fatal("Error:", err)
	}

	testUserRepo = UserModel{DB: testDB}
	testUrlRepo = UrlModel{DB: testDB}
	testTokenRepo = TokenModel{DB: testDB}

	// run tests
	code := m.Run()

	os.Exit(code)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
