package models

import (
	"errors"
	"log"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	if db.DB == nil {
		return errors.New("database connection is nil")
	}
	query := `INSERT INTO users (email, password) VALUES (?,?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.Hash(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted id: %v", err)
		return err
	}
	u.ID = id
	return err
}

//UserID

func (u *User) ValidateCredentials() error {
	if db.DB == nil {
		return errors.New("database connection is nil")
	}
	query := "SELECT id, password FROM users WHERE email =?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(u.Password, retrievedPassword) {
		return errors.New("invalid credentials")
	}
	return nil
}
