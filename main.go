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
		"github.com/joho/godotenv"
		_ "modernc.org/sqlite"
	)

	func connectDb() *sql.DB {
		db, err := sql.Open("sqlite", "./example.db")
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			log.Fatal(err)
		}

		err = godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		return db
	}

	func main() {
		db := connectDb()
		defer db.Close()

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
		c.IndentedJSON(http.StatusOK, quotes)
	}

	func addQuote(c *gin.Context, quotesDao quotes.QuotesDao) {
		// Retrieve the JWT token from the "Authorization" header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			return
		}
		// Parse and validate the JWT token
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
	
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
	
		// Extract user data from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
	
		userIdFloat, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}
		userId := int(userIdFloat)
	
		// Simulate fetching user object (or skip if not needed)
		//authenticatedUser := users.User{Id: userId}
	
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

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})

	}
