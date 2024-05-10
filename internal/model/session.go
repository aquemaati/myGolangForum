package model

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// Session represents a user session
type Session struct {
	SessionID string
	UserID    string
	ExpiresAt time.Time
}

// CreateSession creates a new session and stores it in the database
func CreateSession(db *sql.DB, userID string, duration time.Duration) (Session, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return Session{}, err
	}
	expiresAt := time.Now().Add(duration)

	query := "INSERT INTO Sessions (sessionId, userId, expiresAt) VALUES (?, ?, ?)"
	_, err = ExecuteNonQuery(db, query, sessionID, userID, expiresAt)
	if err != nil {
		return Session{}, err
	}
	return Session{SessionID: sessionID.String(), UserID: userID, ExpiresAt: expiresAt}, nil
}

// GetSession retrieves a session by its ID
func GetSession(db *sql.DB, sessionID string) (Session, error) {
	var s Session
	err := db.QueryRow("SELECT sessionId, userId, expiresAt FROM Sessions WHERE sessionId = ?", sessionID).Scan(&s.SessionID, &s.UserID, &s.ExpiresAt)
	if err != nil {
		return Session{}, err
	}
	return s, nil
}

// SessionValid checks if a session is valid based on the expiry time
func SessionValid(s Session) bool {
	return s.ExpiresAt.After(time.Now())
}

// DeleteSession removes a session from the database
func DeleteSession(db *sql.DB, sessionID string) error {
	query := "DELETE FROM Sessions WHERE sessionId = ?"
	_, err := ExecuteNonQuery(db, query, sessionID)
	return err
}

// CheckSessionExists checks if there is an active session for the given user ID
func CheckSessionExists(db *sql.DB, userID string) bool {
	var expiresAt time.Time
	var sessionId string

	// Prepare and execute the query
	query := "SELECT sessionId, expiresAt FROM Sessions WHERE userId = ?"
	err := db.QueryRow(query, userID).Scan(&sessionId, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false // No active session found for this user
		}
		return false
	}

	// Check if the session has expired
	if time.Now().After(expiresAt) {
		_ = DeleteSession(db, sessionId) // Remove expired session
		return false
	}
	return true
}

func InvalidateSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Unix(0, 0), // Expire in the past
		Path:    "/",
	}
	http.SetCookie(w, cookie)
}
