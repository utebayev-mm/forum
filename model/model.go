package model

import "database/sql"

type User struct {
	Id                 int
	Name               string
	Password           string
	Email              string
	Permissions        string
	Role_id            int
	RequestedPromotion bool
}

type Session struct {
	CookieName  string
	CookieValue string
	UserId      int
}

type Category struct {
	Id          int
	Name        string
	Description string
}

type Report struct {
	Id         int
	UserID     int
	ReportID   int
	Status     string
	Title      string
	Sender     string
	AdminReply string
}

type Post struct {
	Id           int
	Title        string
	Body         string
	Author       int
	AuthorName   string
	Category     int
	CategoryName string
	Tags         string
	SeparateTags []string
	PostingTime  string
	Likes        int
	Dislikes     int
	Image        string
}

type Comment struct {
	CommentID       int
	CommentText     string
	PostID          int
	PostTitle       string
	AuthorID        int
	CommentLikes    int
	CommentDislikes int
	AuthorName      string
	PostingTime     string
	Likes           int
	Dislikes        int
}

type UserLike struct {
	LikeID     int
	Mark       string
	AuthorID   int
	PostID     int
	AuthorName string
}

type CommentUserLike struct {
	CommentLikeID int
	CommentMark   string
	CommentID     int
	UserID        int
	AuthorName    string
}

var PostFeed struct {
	Posts []Post
}

var DB *sql.DB
