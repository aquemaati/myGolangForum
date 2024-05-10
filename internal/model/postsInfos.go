package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type PostInfo struct {
	PostId        int
	UserId        string
	UserImage     string
	UserName      string
	PostDate      time.Time
	LoveNumb      int
	HateNumb      int
	Title         string
	Description   string
	CategoryNames []string
	Comments      []CommentInfo
}

func ScanPostInfo(rows *sql.Rows) (PostInfo, error) {
	var categoryNames string
	var p PostInfo
	err := rows.Scan(
		&p.PostId,
		&p.UserId,
		&p.UserImage,
		&p.UserName,
		&p.PostDate,
		&p.LoveNumb,
		&p.HateNumb,
		&p.Title,
		&p.Description,
		&categoryNames,
	)
	if err != nil {
		return PostInfo{}, err
	}
	// Diviser la chaîne en un tableau
	if categoryNames != "" {
		p.CategoryNames = strings.Split(categoryNames, ",")
	}
	return p, nil
}

func FetchPosts(db *sql.DB, userId *string, category *string) ([]PostInfo, error) {
	query := `
	    SELECT postId, userId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames
	    FROM ExtendedPostsInfosView
	    WHERE 1 = 1`
	args := []interface{}{}

	if userId != nil {
		query += " AND userId = ?"
		args = append(args, *userId)
	}

	if category != nil {
		catPattern := fmt.Sprintf("%%%s%%", *category)
		query += " AND categoryNames LIKE ?"
		args = append(args, catPattern)
	}

	return ExecuteQuery(db, query, ScanPostInfo, args...)
}

func FetchExtendedPostsWithComments(db *sql.DB, userId *string, category *string) ([]PostInfo, error) {
	posts, err := FetchPosts(db, userId, category)
	if err != nil {
		return nil, err
	}

	comments, err := FetchComments(db)
	if err != nil {
		return nil, err
	}

	postMap := make(map[int]*PostInfo)
	for i := range posts {
		postMap[posts[i].PostId] = &posts[i]
	}

	for _, comment := range comments {
		if post, exists := postMap[comment.PostId]; exists {
			post.Comments = append(post.Comments, comment)
		}
	}

	return posts, nil
}

func FetchUniquePost(db *sql.DB, postId int) (PostInfo, error) {
	var postInfo PostInfo
	var categoryNames string

	// Define the query to fetch the specific post by postId
	query := `
		SELECT postId, userId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames
		FROM ExtendedPostsInfosView
		WHERE postId = ?
	`
	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return PostInfo{}, err
	}
	defer stmt.Close()

	// Execute the query and scan the result into postInfo
	err = stmt.QueryRow(postId).Scan(
		&postInfo.PostId,
		&postInfo.UserId,
		&postInfo.UserImage,
		&postInfo.UserName,
		&postInfo.PostDate,
		&postInfo.LoveNumb,
		&postInfo.HateNumb,
		&postInfo.Title,
		&postInfo.Description,
		&categoryNames,
	)
	if err != nil {
		return PostInfo{}, err
	}

	// Convert category names to a slice
	if categoryNames != "" {
		postInfo.CategoryNames = strings.Split(categoryNames, ",")
	}

	// Fetch comments associated with this post
	comments, err := FetchCommentsByPostId(db, postId)
	if err != nil {
		return PostInfo{}, err
	}

	// Assign comments to the post's Comments field
	postInfo.Comments = comments

	return postInfo, nil
}
