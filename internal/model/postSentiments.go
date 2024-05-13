package model

import (
	"database/sql"
	"errors"
)

var ErrSentimentNotFound = errors.New("sentiment not found")

// AddPostLike adds or updates a like or dislike for a post by a user
func AddPostLike(db *sql.DB, userID string, postId string, sentiment string) error {
	if sentiment != "love" && sentiment != "hate" {
		return errors.New("invalid sentiment")
	}

	query := `
        INSERT INTO PostsLike (userId, postId, sentiment)
        VALUES (?, ?, ?)
        ON CONFLICT (userId, postId) DO UPDATE SET sentiment = excluded.sentiment
    `
	_, err := ExecuteNonQuery(db, query, userID, postId, sentiment)
	return err
}

func GetUserSentiment(db *sql.DB, userId string, postId string) (string, error) {
	var sentiment string
	query := `SELECT sentiment FROM PostsLike WHERE userId = ? AND postId = ?`
	err := db.QueryRow(query, userId, postId).Scan(&sentiment)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrSentimentNotFound
		}
		return "", err
	}
	return sentiment, nil
}

// RemovePostLike removes a sentiment associated with a user and a post
func RemovePostLike(db *sql.DB, userID string, postId string) error {
	query := `DELETE FROM PostsLike WHERE userId = ? AND postId = ?`
	_, err := ExecuteNonQuery(db, query, userID, postId)
	return err
}
