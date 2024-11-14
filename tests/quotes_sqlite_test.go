package quotes_test

import (
	"database/sql"
	"testing"
	"time"
	"whattoday/web-service-gin/quotes"
	"whattoday/web-service-gin/users"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("Failed to enable foreign key constraints: %v", err)
	}

	quotesSqlite := quotes.NewQuotesSqlite(db)
	quotesSqlite.CreateTable()

	usersSqlite := users.NewUserSqlite(db)
	usersSqlite.CreateTable()

	return db
}

func TestCreateQuotesTable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	quotesSqlite := quotes.NewQuotesSqlite(db)
	assert.NotPanics(t, func() { quotesSqlite.CreateTable() }, "CreateTable should not panic")
}

func TestAddQuote(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	quotesSqlite := quotes.NewQuotesSqlite(db)

	_, err := db.Exec(`INSERT INTO users (id, name) VALUES (1, 'Test User')`)
	assert.NoError(t, err)

	quote := quotes.Quote{
		Text:        "This is a test quote",
		UserId:      1,
		PublishDate: time.Now(),
	}
	err = quotesSqlite.AddQuote(quote)
	assert.NoError(t, err)

}

func TestAddQuoteFailsForeingKey(t *testing.T){
	db := setupTestDB(t)
	defer db.Close()

	quotesSqlite := quotes.NewQuotesSqlite(db)


	_, err := db.Exec(`INSERT INTO users (id, name) VALUES (1, 'Test User')`)
	assert.NoError(t, err)

	invalidQuote := quotes.Quote{
		Text:        "Invalid quote",
		UserId:      999,
		PublishDate: time.Now(),
	}
	err = quotesSqlite.AddQuote(invalidQuote)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FOREIGN KEY constraint failed")

}

func TestGetQuoteThatExists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	quotesSqlite := quotes.NewQuotesSqlite(db)

	_, err := db.Exec(`INSERT INTO users (id, name) VALUES (1, 'Test User')`)
	assert.NoError(t, err)


	quote := quotes.Quote{
		Text:        "Sample quote",
		UserId:      1,
		PublishDate: time.Now(),
	}
	err = quotesSqlite.AddQuote(quote)
	assert.NoError(t, err)


	quotesList, err := quotesSqlite.GetQuotes()
	assert.NoError(t, err)
	assert.Len(t, quotesList, 1)
	assert.Equal(t, "Sample quote", quotesList[0].Text)
	assert.Equal(t, 1, quotesList[0].UserId)
}

