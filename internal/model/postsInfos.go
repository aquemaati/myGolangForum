package model

import "database/sql"

type PostInfo struct {
	PostId        int
	UserImage     string
	UserName      string
	PostDate      string
	LoveNumb      int
	HateNumb      int
	Title         string
	Description   string
	CategoryNames string
	Comments      []CommentInfo
}

type CommentInfo struct {
	PostId    int
	CommentId int
	UserImage string
	UserName  string
	Date      string
	LoveNumb  int
	HateNumb  int
	Content   string
}

func ScanPostInfo(rows *sql.Rows) (PostInfo, error) {
	var p PostInfo
	err := rows.Scan(
		&p.PostId,
		&p.UserImage,
		&p.UserName,
		&p.PostDate,
		&p.LoveNumb,
		&p.HateNumb,
		&p.Title,
		&p.Description,
		&p.CategoryNames,
	)
	if err != nil {
		return PostInfo{}, err
	}
	return p, nil
}

func ScanCommentInfo(rows *sql.Rows) (CommentInfo, error) {
	var c CommentInfo
	err := rows.Scan(
		&c.PostId,
		&c.CommentId,
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

func FetchExtendedPostsWithComments(db *sql.DB) ([]PostInfo, error) {
	// Récupérer les posts
	queryPosts := "SELECT postId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames FROM ExtendedPostsInfosView"
	posts, err := ExecuteQuery(db, queryPosts, ScanPostInfo)
	if err != nil {
		return nil, err
	}

	// Récupérer les commentaires
	queryComments := "SELECT postId, commentId, userImage, userName, date, loveNumb, hateNumb, content FROM CommentsInfosView"
	comments, err := ExecuteQuery(db, queryComments, ScanCommentInfo)
	if err != nil {
		return nil, err
	}

	// Associer les commentaires aux posts
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
