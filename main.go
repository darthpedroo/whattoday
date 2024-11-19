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

		"github.com/gin-contrib/cors"
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
		config := cors.DefaultConfig()
		
		config.AllowOrigins = []string{
			"http://localhost:5500",  
			"https://darthpedroo.github.io",  
		}
		config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
		config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
		config.ExposeHeaders = []string{"Content-Length"}
		config.AllowCredentials = true
		config.MaxAge = 12 * time.Hour

		router.Use(cors.New(config))

		router.GET("/quotes", func(c *gin.Context) {
			getQuotes(c, quotesSqlite)
		})
		router.POST("/quotes", func(c *gin.Context) {
			middleware.RequireAuth(usersSqlite)(c) // Call RequireAuth with userDao, passing c
			addQuote(c, quotesSqlite)
		})

		router.POST("/sign-up", func(c *gin.Context) {
			addUser(c, usersSqlite)
		})
		router.POST("/login", func(c *gin.Context) {
			Login(c, usersSqlite)
		})

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080" // Fallback en caso de que no esté definida
		}
		router.Run(":" + port)
		
		}

	func getQuotes(c *gin.Context, quotesDao quotes.QuotesDao) {
		quotes, err := quotesDao.GetQuotes()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Print(quotes)
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
		currentUserFromDb, err := userDao.GetUserFromName(newUser.Name)

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
			"exp":  time.Now().Add(time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.SetSameSite(http.SameSiteNoneMode) // Allows cross-origin cookies
	
		environment := os.Getenv("APP_ENV") // Get the environment from the environment variable

		if environment == "production" {
			fmt.Println("xd?")
			c.SetSameSite(http.SameSiteNoneMode) // Required for cross-origin cookies
			c.SetCookie(
				"Authorization",   // Cookie name
				tokenString,       // Cookie value
				3600,              // Expiry time in seconds
				"/",               // Path
				".github.io",      // Domain
				true,              // Secure
				false,             // HttpOnly
			)
		} else {
			fmt.Println("holii")
			c.SetSameSite(http.SameSiteLaxMode) // More relaxed in development
			c.SetCookie(
				"Authorization",   // Cookie name
				tokenString,       // Cookie value
				3600,              // Expiry time in seconds
				"/",               // Path
				"localhost",       // Domain for local development
				false,             // Non-secure in development (no HTTPS)
				false,              // HttpOnly
			)
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
		fmt.Println("set the cookie :V")

	}
