package model

import "time"

type CommentInfos struct {
	PostId    int
	CommentId int
	UserImage string
	UserName  string
	Date      time.Time
	LoveNumb  int
	HateNumb  int
	Content   string
}
