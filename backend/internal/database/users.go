package database

import (
	"database/sql"
	"errors"
	"fileserverapi/internal/auth"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username string, password string) {
	hash, err := hashPassword(password)
	if err != nil {
		log.Println("Could not hash password", err)
	}

	_, err = DB.Query("INSERT INTO users (username,password) VALUES ($1,$2)", username, hash)
	if err != nil {
		log.Println("Could not insert user", err)
	}
}

func Authenticate(username string, password string) (string, error) {
	var pw string

	err := DB.QueryRow("SELECT users.password FROM users WHERE users.username = $1", username).Scan(&pw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("401 user notz found")
			return "", err
		}
		log.Println("Could not exec the query", err)
	}

	if !auth.CheckPasswordHash(pw, password) {
		return "", errors.New("Could not check password")
	}

	t, err := auth.GenerateJWT(username)
	if err != nil {
		log.Printf("Could not generate jwt: %s", err)
		return "", nil
	}

	return t, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}
