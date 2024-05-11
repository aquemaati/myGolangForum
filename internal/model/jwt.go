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
func GenerateJWT(db *sql.DB, userId string, contextKey string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	encodedHeader := encodeBase64(headerJSON)

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := map[string]interface{}{
		contextKey: userId,
		"exp":      expirationTime.Unix(),
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

func ParseJWT(token string) (map[string]interface{}, error) {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payloadSegment := segments[1]
	payloadData, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(payloadSegment)
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	err = json.Unmarshal(payloadData, &claims)
	if err != nil {
		return nil, err
	}

	signatureSegment := segments[2]
	expectedSignature := signMessage(segments[0] + "." + segments[1])
	if signatureSegment != expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}
	} else {
		return nil, fmt.Errorf("expiration claim is missing")
	}

	return claims, nil
}
