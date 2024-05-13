package model

import (
	"database/sql"
)

type CommentInfo struct {
	PostId    string
	CommentId string
	UserId    string
	UserImage string
	UserName  string
	Date      string
	LoveNumb  int
	HateNumb  int
	Content   string
}

func ScanCommentInfo(rows *sql.Rows) (CommentInfo, error) {
	var c CommentInfo
	err := rows.Scan(
		&c.PostId,
		&c.CommentId,
		&c.UserId,
		&c.UserImage,
		&c.UserName,
		&c.Date,
		&c.LoveNumb,
		&c.HateNumb,
		&c.Content,
	)
	if err != nil {
		return CommentInfo{}, err
	}
	return c, nil
}

func FetchComments(db *sql.DB) ([]CommentInfo, error) {
	query := `
	    SELECT postId, commentId, userId, userImage, userName, date, loveNumb, hateNumb, content
	    FROM CommentsInfosView`
	return ExecuteQuery(db, query, ScanCommentInfo)
}

func FetchCommentsByPostId(db *sql.DB, postId string) ([]CommentInfo, error) {
	query := `
	    SELECT postId, commentId, userId, userImage, userName, date, loveNumb, hateNumb, content
	    FROM CommentsInfosView WHERE postId = ?`
	return ExecuteQuery(db, query, ScanCommentInfo, postId)
}
