package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"whattoday/web-service-gin/middleware"
	"whattoday/web-service-gin/quotes"
	"whattoday/web-service-gin/users"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite", "./example.db")
	_, err = db.Exec("PRAGMA foreign_keys = ON")
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
	router.POST("/quotes", func(c *gin.Context) {
		middleware.RequireAuth(usersSqlite)(c) // Call RequireAuth with userDao, passing c
		addQuote(c, quotesSqlite)
	})

	router.POST("/users", func(c *gin.Context) {
		addUser(c, usersSqlite)
	})
	router.POST("/login", func(c *gin.Context) {
		Login(c, usersSqlite)
	})

	router.Run("localhost:8080")
}

func getQuotes(c *gin.Context, quotesDao quotes.QuotesDao) {
	quotes, err := quotesDao.GetQuotes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(quotes)
	c.IndentedJSON(http.StatusOK, quotes)
}

func addQuote(c *gin.Context, quotesDao quotes.QuotesDao) {

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	authenticatedUser, ok := user.(users.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User data is incorrect"})
		return
	}

	userId := authenticatedUser.Id
	fmt.Println("Authenticated UserId:", userId)


	var newQuote quotes.Quote
	if err := c.BindJSON(&newQuote); err != nil {
		return
	}

	

	newQuote.PublishDate = time.Now()
	newQuote.UserId = userId

	if err := quotesDao.AddQuote(newQuote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newQuote)
}

func addUser(c *gin.Context, userDao users.UserDao) {

	var newUser users.User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	if err := userDao.AddUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newUser)
}

func Login(c *gin.Context, userDao users.UserDao) {

	// get the user and password of the req body
	var newUser users.User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	// look up requested user
	currentUserFromDb, err := userDao.GetUser(newUser.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Compare sent in pass with saved user pass hash
	fmt.Println("XD")
	fmt.Println(currentUserFromDb.Password)

	err = bcrypt.CompareHashAndPassword([]byte(currentUserFromDb.Password), []byte(newUser.Password))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate  a jwt token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  currentUserFromDb.Id,
		"name": currentUserFromDb.Name,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}
