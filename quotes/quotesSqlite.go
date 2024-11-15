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

func (q QuotesSqlite) GetQuotes() ([]CustomQuoteResponse, error) {

	

	db := q.db

	//var quotes = make([]Quote, 0)
	//var usersList = make([]users.User, 0)

	rows, err := db.Query(`select q.id,q.text,q.publishDate,q.userId, u.id, u.name from quotes q 
						join users u on q.userId = u.id `)

	var response []CustomQuoteResponse

	for rows.Next() {
		
		var quote Quote
		var user users.User

		// se escanea el espacio de memoria pq la respuesta es algo asi: // Return error if scanning fails
		if err := rows.Scan(&quote.Id, &quote.Text, &quote.PublishDate, &quote.UserId, &user.Id, &user.Name); err != nil {
			log.Printf("Error scanning quote: %v", err)
			return nil, err
		}
		
		customResponse := CustomQuoteResponse{
			Quote: quote,
			User:  user,
		}

		response = append(response, customResponse)
		
		//quotes = append(quotes, quote)
		//usersList = append(usersList, user)
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return response, nil
}

func (q QuotesSqlite) AddQuote(quote Quote) error {
	db := q.db

	canPost, hoursUntilNextPost, err := users.UserCanPost(db, quote.UserId)

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
