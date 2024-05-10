package model

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	UserName string
	Email    string
	Password string
	Date     time.Time
	Image    string
}

type UserPublic struct {
	UserName string
	Image    string
}

func ScanUserInfo(row *sql.Row) (UserPublic, error) {
	var user UserPublic
	err := row.Scan(&user.UserName, &user.Image)
	if err != nil {
		return UserPublic{}, err
	}
	return user, nil
}

func CreateUserInDB(db *sql.DB, username, email, password string) (User, error) {
	// Generate a new UUID for the user
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return User{}, err
	}

	//Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return User{}, err
	}

	query := "INSERT INTO Users (id, userName, email, password, image) VALUES (?, ?, ?, ?, ?)"
	_, err = ExecuteNonQuery(db, query, id.String(), username, email, string(hashedPassword), "path/default/image.jpg")
	if err != nil {
		return User{}, err
	}

	return User{
		ID:       id.String(),
		UserName: username,
	}, nil
}

func GetUserById(db *sql.DB, userId string) (UserPublic, error) {
	query := `SELECT userName, image FROM Users WHERE id = ?`
	return ExecuteSingleQuery(db, query, ScanUserInfo, userId)
}
