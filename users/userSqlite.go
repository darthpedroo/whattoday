package users

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite(db *sql.DB) *UserSqlite {
	return &UserSqlite{db: db}
}

func (u UserSqlite) GetUsers() ([]User, error) {
	db := u.db

	var users = make([]User, 0)

	rows, err := db.Query(`SELECT * FROM users`)

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			log.Printf("Error scanning User: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return users, nil
}

func (u UserSqlite) AddUser(user User) error {

	db := u.db
	stmt, err := db.Prepare(`INSERT INTO users (name,password)  VALUES (?,?) `)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	_, err = stmt.Exec(user.Name, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (u UserSqlite) GetUser(userId int) (User, error) {

	db := u.db
	row, err := db.Query(`SELECT * FROM users where id = ?`, userId)

	var user User
	for row.Next() {

		if err := row.Scan(&user.Id, &user.Name, &user.Password); err != nil {
			log.Printf("Error scanning User: %v", err)
			return User{}, err
		}
	}

	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	return user, nil

}

func (u UserSqlite) GetUserFromName(userName string) (User, error) {

	db := u.db
	row, err := db.Query(`SELECT * FROM users where name = ?`, userName)

	var user User
	for row.Next() {

		if err := row.Scan(&user.Id, &user.Name, &user.Password); err != nil {
			log.Printf("Error scanning User: %v", err)
			return User{}, err
		}
	}

	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	return user, nil

}

func UserCanPost(db *sql.DB, userId int) (bool, time.Duration, error) {

	row, err := db.Query(`SELECT q.publishDate FROM users u
								JOIN quotes q ON u.id = q.userId
								WHERE q.userId = ?
								order by q.publishDate DESC
								limit 1;
								
							`, userId)

	if err != nil {
		log.Fatal(err)
		return false, 0, err
	}
	var publishDate time.Time //This stores the last published date
	for row.Next() {

		err := row.Scan(&publishDate)

		if err != nil {
			return false, 0, err
		}

	}

	currentTime := time.Now()

	expectedTime := publishDate.AddDate(0, 0, 1)

	hoursUntilNextPost := expectedTime.Sub(currentTime)

	canPost := hoursUntilNextPost < 0

	defer row.Close()
	return canPost, hoursUntilNextPost, nil

}

func (u UserSqlite) CreateTable() {
	db := u.db

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	password TEXT)`)

	if err != nil {
		log.Fatal(err)
	}
}
