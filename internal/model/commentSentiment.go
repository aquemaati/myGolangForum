package model

import (
	"database/sql"
	"errors"
)

func AddSentiimentComment(db *sql.DB, userID string, commentId string, sentiment string) error {
	if sentiment != "love" && sentiment != "hate" {
		return errors.New("invalid sentiment")
	}

	query := `
	    INSERT INTO CommentsLike (userId, commentsId, sentiment)
	    VALUES (?, ?, ?)
	    ON CONFLICT (userId, commentsId) DO UPDATE SET sentiment = excluded.sentiment
	`

	_, err := ExecuteNonQuery(db, query, userID, commentId, sentiment)
	return err
}

func GetUserSentimentComment(db *sql.DB, userId string, postId string) (string, error) {
	query := `SELECT sentiment FROM CommentsLike WHERE userId = ? AND commentsId = ?`

	var sentiment string
	err := db.QueryRow(query, userId, postId).Scan(&sentiment)

	if err == sql.ErrNoRows {
		return "", ErrSentimentNotFound
	}

	if err != nil {
		return "", err
	}

	return sentiment, nil
}

func RemoveSentomentComment(db *sql.DB, userID string, postId string) error {
	query := `DELETE FROM CommentsLike WHERE userId = ? AND commentsId = ?`
	_, err := ExecuteNonQuery(db, query, userID, postId)
	return err
}
