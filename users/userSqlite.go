package users

import (
	"database/sql"
	"log"
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
	stmt, err := db.Prepare(`INSERT INTO users (name) VALUES (?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (u UserSqlite) CreateTable() {
	db := u.db

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT)`)

	if err != nil {
		log.Fatal(err)
	}
}
