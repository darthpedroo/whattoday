package quotes

import (
	"database/sql"
	"fmt"
	"log"
	"whattoday/web-service-gin/users"
)

type QuotesSqlite struct {
	db *sql.DB
}

func NewQuotesSqlite(db *sql.DB) *QuotesSqlite {
	return &QuotesSqlite{db: db}
}

func (q QuotesSqlite) GetQuotes() ([]Quote, error) {
	db := q.db

	var quotes = make([]Quote, 0)

	rows, err := db.Query(`SELECT * FROM quotes`)

	for rows.Next() {
		var quote Quote
		// se escanea el espacio de memoria pq la respuesta es algo asi: // Return error if scanning fails
		if err := rows.Scan(&quote.Id, &quote.Text, &quote.PublishDate, &quote.UserId); err != nil {
			log.Printf("Error scanning quote: %v", err)
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return quotes, nil

}

func (q QuotesSqlite) AddQuote(quote Quote) error {
	db := q.db

	canPost, hoursUntilNextPost, err := users.UserCanPost(db, quote.UserId)
	
	fmt.Println("CAN POST: ")
	fmt.Println(canPost)
	fmt.Println("HOURS UNTIL NEXT POST")
	fmt.Println(hoursUntilNextPost)

	if !canPost {
		return fmt.Errorf("user cannot post yet, hours till next post: %v", hoursUntilNextPost)
	}

	stmt, err := db.Prepare(`INSERT INTO quotes (text, userId, publishDate)  VALUES (?,?,?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(quote.Text, quote.UserId, quote.PublishDate)
	if err != nil {
		return err
	}

	return nil

}

func (q QuotesSqlite) CreateTable() {
	db := q.db

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		publishDate DATETIME,
		userId INTEGER,
		FOREIGN KEY (userId) REFERENCES users(id)
	)`)
	if err != nil {
		log.Fatal(err)
	}

}
