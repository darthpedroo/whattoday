package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"whattoday/web-service-gin/quotes"
	"whattoday/web-service-gin/users"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite", "./example.db")
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to SQLite database!")
	return db
}

func main() {

	db := connectDb()
	defer db.Close()
	fmt.Println(db)

	quotesSqlite := quotes.NewQuotesSqlite(db)
	quotesSqlite.CreateTable()
	usersSqlite := users.NewUserSqlite(db)
	usersSqlite.CreateTable()

	router := gin.Default()
	router.GET("/quotes", func(c *gin.Context) {
		getQuotes(c, quotesSqlite)
	})
	router.POST("quotes", func(c *gin.Context) {
		addQuote(c, quotesSqlite)
	})
	router.POST("users", func(c *gin.Context){

	})

	router.Run("localhost:8080")
}

//func addUser(c *gin.Context)

func getQuotes(c *gin.Context, quotesDao quotes.QuotesDao) {
	quotes, err := quotesDao.GetQuotes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(quotes)
	c.IndentedJSON(http.StatusOK, quotes)
}

func addQuote(c *gin.Context, quotesDao quotes.QuotesDao) {
	
	var newQuote quotes.Quote
	if err := c.BindJSON(&newQuote); err != nil {
		return
	}

	newQuote.PublishDate = time.Now()

	quotesDao.AddQuote(newQuote)
	c.IndentedJSON(http.StatusCreated, newQuote)
}
