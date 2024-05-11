package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var SecretKey = []byte("your_secret_key")

func encodeBase64(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func signMessage(message string) string {
	mac := hmac.New(sha256.New, SecretKey)
	mac.Write([]byte(message))
	return encodeBase64(mac.Sum(nil))
}

// function to generate a jwt
func GenerateJWT(db *sql.DB, userId string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	encodedHeader := encodeBase64(headerJSON)

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := map[string]interface{}{
		"userId": userId,
		"exp":    expirationTime.Unix(),
	}
	claimsJSON, _ := json.Marshal(claims)
	encodedClaims := encodeBase64(claimsJSON)

	message := fmt.Sprintf("%s.%s", encodedHeader, encodedClaims)
	signature := signMessage(message)
	jwt := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedClaims, signature)

	// Generate a UUID for the session ID
	sessionID := uuid.New().String()

	// Insert the session into the database TODO adapt with util foction
	_, err := db.Exec("INSERT INTO Sessions (SessionID, UserID, JWT, ExpiresAt) VALUES (?, ?, ?, ?)",
		sessionID, userId, jwt, expirationTime)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
