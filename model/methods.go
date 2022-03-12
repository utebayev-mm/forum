package model

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func (u User) Create() (int, error) {
	insertUserSQL := `INSERT INTO user(name, password, email, role_id) VALUES (?, ?, ?, ?)`
	statement, err := DB.Prepare(insertUserSQL)

	if err != nil {
		return 0, err
	}
	res, err := statement.Exec(u.Name, u.Password, u.Email, u.Role_id)
	statement.Close()
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err

	}
	return int(id), nil
}

func (c Comment) Create() error {
	insertUserSQL := `INSERT INTO comment(content, user_id, post_id, posting_time) VALUES (?, ?, ?, ?)`
	statement, err := DB.Prepare(insertUserSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(c.CommentText, c.AuthorID, c.PostID, c.PostingTime)
	statement.Close()
	if err != nil {
		return err

	}
	return nil
}

func DeleteComment(id int) error {
	sqlStatement := `
	DELETE FROM comment
	WHERE id = $1;`
	statement, err := DB.Prepare(sqlStatement)

	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	statement.Close()
	if err != nil {
		return err

	}
	return nil
}

func (l UserLike) Delete() error {
	sqlStatement := `
DELETE FROM userlike
WHERE user_id = $1 AND post_id = $2;`
	_, err := DB.Exec(sqlStatement, l.AuthorID, l.PostID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategory(id int) error {
	sqlStatement := `
DELETE FROM category
WHERE id = $1;`
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func FindAndDeletePostsByCategory(category_id int) {
	var posts []int
	rows, err := DB.Query("select id from post WHERE category_id = ?", category_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		posts = append(posts, id)
	}
	fmt.Println(posts)
	for _, value := range posts {
		DeleteCommentInMyComments(value)
		DeletePostInActivity(value)
		fmt.Println("post to delete id=", value)
		err := DeletePost(value)
		if err != nil {
			log.Println(err)
		}
		err = DeleteLikeAfterPost(value)
		if err != nil {
			log.Println(err)
		}
	}
}
func DeletePostsFromCategory(id int) error {
	sqlStatement := `
DELETE FROM post
WHERE category_id = $1;`
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func (l CommentUserLike) Delete() error {
	sqlStatement := `
DELETE FROM userlike_comment
WHERE user_id = $1 AND comment_id = $2;`
	_, err := DB.Exec(sqlStatement, l.UserID, l.CommentID)
	if err != nil {
		return err
	}
	return nil
}

func (l UserLike) Create() error {
	insertUserSQL := `INSERT INTO userlike(mark, user_id, post_id) VALUES (?, ?, ?)`
	statement, err := DB.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	if l.PostID != 0 {
		_, err = statement.Exec(l.Mark, l.AuthorID, l.PostID)
	}
	statement.Close()
	if err != nil {
		return err
	}
	return nil
}

func DeleteLikeAfterPost(post_id int) error {
	sqlStatement := `
DELETE FROM userlike
WHERE post_id = $1;`
	_, err := DB.Exec(sqlStatement, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (u User) SetUserIdByEmail() error {
	var id string
	err := DB.QueryRow("select id from user where email=?", u.Email).Scan(&id)
	if err != nil {
		return err
	}
	uId, _ := strconv.Atoi(id)
	u.Id = uId
	return nil
}

// is user need cookie
func (u User) IsUserNeedCookie() bool {
	var id string
	err := DB.QueryRow("select id from session where id=?", u.Id).Scan(&id)

	if err != nil {
		return true
	}
	return false
}

func (u User) UserIdByEmail() int {
	var id string
	DB.QueryRow("select id from user where email=?", u.Email).Scan(&id)
	uId, _ := strconv.Atoi(id)
	return uId
}

func (s Session) Create() error {
	insertUserSQL := `INSERT INTO session(cookiename, cookievalue, user_id) VALUES (?, ?, ?)`
	statement, err := DB.Prepare(insertUserSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(s.CookieName, s.CookieValue, s.UserId)
	statement.Close()

	if err != nil {
		return err

	}
	return nil
}

func (u User) IsPasswordCorrect() bool {
	var passwordHash string
	err := DB.QueryRow("select password from user where email=?", u.Email).Scan(&passwordHash)
	if err != nil {
		log.Println(err)
		return false
	}

	if !CheckPassword(passwordHash, u.Password) {
		return false
	}
	return true
}

func (p Post) Create() error {
	insertUserSQL := `INSERT INTO post(title, body, user_id, tags, image, category_id, posting_time) VALUES (?, ?, ?, ?, ?, ?, ?)`
	statement, err := DB.Prepare(insertUserSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(p.Title, p.Body, p.Author, p.Tags, p.Image, p.Category, p.PostingTime)
	statement.Close()

	if err != nil {
		return err

	}
	return nil
}

func (p Post) Update(id int) error {
	updatePostSQL := `UPDATE post SET title=$1, body=$2, user_id=$3, tags=$4, image=$5, category_id=$6, posting_time=$7 WHERE id=$8;`
	statement, err := DB.Prepare(updatePostSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(p.Title, p.Body, p.Author, p.Tags, p.Image, p.Category, p.PostingTime, id)
	statement.Close()

	if err != nil {
		return err

	}
	return nil
}

func UpdateUserRole(id, role_id int) error {
	updateUserSQL := `UPDATE user SET role_id=$1 WHERE id=$2;`
	statement, err := DB.Prepare(updateUserSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(role_id, id)
	statement.Close()

	if err != nil {
		return err

	}
	return nil
}

func DeleteRequest(user_id int) {
	updateUserSQL := `DELETE FROM requests WHERE user_id=$1;`
	statement, err := DB.Prepare(updateUserSQL)

	if err != nil {
		fmt.Println(err)
	}
	_, err = statement.Exec(user_id)
	statement.Close()

	if err != nil {
		fmt.Println(err)
	}
}

func (c Comment) Update(id int) error {
	updateCommentSQL := `UPDATE comment SET content=$1 WHERE id=$2;`
	statement, err := DB.Prepare(updateCommentSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(c.CommentText, id)
	statement.Close()

	if err != nil {
		return err

	}
	return nil
}

func DeleteReport(id int) error {
	updateReportSQL := `DELETE FROM reports WHERE post_id=$2;`
	statement, err := DB.Prepare(updateReportSQL)

	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	statement.Close()

	if err != nil {
		return err

	}
	return nil
}

// func DeletePost(id, userId int) error {
// 	DeleteImage(id)
// 	sqlStatement := `
// DELETE FROM post
// WHERE id = $1 and user_id = $2;`
// 	_, err := DB.Exec(sqlStatement, id, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func DeletePost(id int) error {
	DeleteImage(id)
	sqlStatement := `
DELETE FROM post
WHERE id = $1;`
	_, err := DB.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func DeletePostInActivity(id int) error {
	sqlStatement := `
DELETE FROM user_activity
WHERE post_id = $1;`
	_, err := DB.Exec(sqlStatement, id, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeletePostInReports(id int) error {
	sqlStatement := `
DELETE FROM reports
WHERE post_id = $1;`
	_, err := DB.Exec(sqlStatement, id, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentInMyComments(id int) error {
	fmt.Println("comments from the post ", id)
	sqlStatement := `
		DELETE FROM comment
		WHERE post_id = $1;`
	statement, err := DB.Prepare(sqlStatement)

	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	statement.Close()
	if err != nil {
		return err

	}
	return nil
}

func DeleteImage(postID int) error {
	var image string
	err := DB.QueryRow("select image from post where id=?", postID).Scan(&image)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("FILE TO DELETE", image)
	fileToDeleteErr := os.Remove(image)
	if fileToDeleteErr != nil {
		fmt.Println(fileToDeleteErr)
	}
	fmt.Println("Image from postID ", postID, "deleted")
	return nil
}

func GetAllPosts() []Post {
	var posts []Post
	rows, err := DB.Query("select id, title, body, user_id, category_id, tags, posting_time, image from post")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// var reverseposts []Post
	for rows.Next() {
		var id, user_id, category_id int
		var title, body, tags, posting_time, image string
		rows.Scan(&id, &title, &body, &user_id, &category_id, &tags, &posting_time, &image)
		realcategoryname, _ := GetCategoryNameById(category_id)
		fmt.Println("user_id", user_id)
		realusername, _ := GetUserName(user_id)
		sepTags := strings.Split(tags, " ")
		tagsMap := make(map[string]bool)
		tagList := []string{}
		for _, entry := range sepTags {
			if _, value := tagsMap[entry]; !value {
				tagsMap[entry] = true
				tagList = append(tagList, entry)
			}
		}
		marks := GetLikes(id)
		var likes, dislikes int
		for i := 0; i < len(marks); i++ {
			if marks[i].Mark == "true" {
				likes++
			} else if marks[i].Mark == "false" {
				dislikes++
			}
		}
		post := Post{
			Id:           id,
			Title:        title,
			Body:         body,
			Author:       user_id,
			AuthorName:   realusername,
			Category:     category_id,
			CategoryName: realcategoryname, //categoryName,
			SeparateTags: tagList,
			PostingTime:  posting_time,
			Likes:        likes,
			Dislikes:     dislikes,
			Image:        image,
		}
		posts = append(posts, post)
		// reverseposts = SortPosts(posts)
		sort.SliceStable(posts, func(i, j int) bool { return posts[i].Id > posts[j].Id })
	}
	return posts
}

func GetAllUsers() []User {
	var users []User
	rows, err := DB.Query("select id, name, email, role_id from user")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// var reverseposts []Post
	for rows.Next() {
		var id, role_id int
		var name, email string
		rows.Scan(&id, &name, &email, &role_id)
		user := User{
			Id:      id,
			Name:    name,
			Email:   email,
			Role_id: role_id,
		}
		users = append(users, user)
	}
	return users
}

func GetAllRequests() []int {
	var users []int
	rows, err := DB.Query("select user_id from requests")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// var reverseposts []Post
	for rows.Next() {
		var user_id int
		rows.Scan(&user_id)
		users = append(users, user_id)
	}
	return users
}

func GetAllCategories() []Category {
	var categories []Category
	rows, err := DB.Query("select id, name from category")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// var reverseposts []Post
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		category := Category{
			Id:   id,
			Name: name,
		}
		categories = append(categories, category)
	}
	return categories
}

func GetAllReports() []Report {
	var reports []Report
	rows, err := DB.Query("select id, post_id, user_id, status, admin_reply from reports")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// var reverseposts []Post
	for rows.Next() {
		var id, user_id, post_id int
		var status, username, posttitle, adminReply string
		var err2 error
		rows.Scan(&id, &post_id, &user_id, &status, &adminReply)
		username, err = GetUserName(user_id)
		posttitle, err2 = GetPostTitle(post_id)
		if err != nil {
			fmt.Println(err)
		}
		if err2 != nil {
			fmt.Println(err2)
		}
		report := Report{
			Id:         id,
			UserID:     user_id,
			ReportID:   post_id,
			Status:     status,
			Title:      posttitle,
			Sender:     username,
			AdminReply: adminReply,
		}
		reports = append(reports, report)
	}
	return reports
}

func GetModeratorReports(user_id int) []Report {
	var reports []Report
	rows, err := DB.Query("select id, post_id, user_id, status, admin_reply from reports WHERE user_id=?", user_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// var reverseposts []Post
	for rows.Next() {
		var id, user_id, post_id int
		var status, username, posttitle, adminReply string
		var err2 error
		rows.Scan(&id, &post_id, &user_id, &status, &adminReply)
		username, err = GetUserName(user_id)
		posttitle, err2 = GetPostTitle(post_id)
		if err != nil {
			fmt.Println(err)
		}
		if err2 != nil {
			fmt.Println(err2)
		}
		report := Report{
			Id:         id,
			UserID:     user_id,
			ReportID:   post_id,
			Status:     status,
			Title:      posttitle,
			Sender:     username,
			AdminReply: adminReply,
		}
		reports = append(reports, report)
	}
	return reports
}

func GetPostToEdit(postID int) Post {
	var id, user_id, category_id int
	var title, body, tags, posting_time, image string
	err := DB.QueryRow("select id, title, body, user_id, category_id, tags, image, posting_time from post where id=?", postID).Scan(&id, &title, &body, &user_id, &category_id, &tags, &image, &posting_time)
	if err != nil {
		log.Println(err)
	}
	realcategoryname, _ := GetCategoryNameById(category_id)
	realusername, _ := GetUserName(user_id)
	sepTags := strings.Split(tags, " ")
	tagsMap := make(map[string]bool)
	tagList := []string{}
	for _, entry := range sepTags {
		if _, value := tagsMap[entry]; !value {
			tagsMap[entry] = true
			tagList = append(tagList, entry)
		}
	}
	fmt.Println("sepTags", sepTags)
	fmt.Println("tagList", tagList)
	// fmt.Println("image", image)
	post := Post{
		Id:           id,
		Title:        title,
		Body:         body,
		Author:       user_id,
		AuthorName:   realusername,
		Category:     category_id,
		CategoryName: realcategoryname, //categoryName,
		Tags:         tags,
		SeparateTags: tagList,
		PostingTime:  posting_time,
		Image:        image,
	}
	return post
}

func GetCommentToEdit(commentID int) Comment {
	var id, user_id, post_id int
	var content, posting_time string
	err := DB.QueryRow("select id, content, user_id, post_id, posting_time from comment where id=?", commentID).Scan(&id, &content, &user_id, &post_id, &posting_time)
	if err != nil {
		log.Println(err)
	}
	realusername, _ := GetUserName(user_id)
	comment := Comment{
		CommentID:   id,
		CommentText: content,
		PostID:      post_id,
		AuthorID:    user_id,
		AuthorName:  realusername,
		PostingTime: posting_time,
	}
	return comment
}

func GetActivityPostingTime(post_id int) string {
	var posting_time string
	err := DB.QueryRow("select posting_time from user_activity where notification_type='COMMENT' AND post_id=?", post_id).Scan(&posting_time)
	if err != nil {
		log.Println(err)
	}
	return posting_time
}

func GetTagPosts(tag string) []Post {
	var tagPosts []Post
	allPosts := GetAllPosts()
	for i := 0; i < len(allPosts); i++ {
		for j := 0; j < len(allPosts[i].SeparateTags); j++ {
			if allPosts[i].SeparateTags[j] == tag {
				tagPosts = append(tagPosts, allPosts[i])
			}
		}
	}
	return tagPosts
}

func GetUserPosts(r *http.Request) []Post {
	userID := GetUserId(r)
	var userPosts []Post
	allPosts := GetAllPosts()
	for i := 0; i < len(allPosts); i++ {
		if allPosts[i].Author == userID {
			userPosts = append(userPosts, allPosts[i])
		}
	}
	return userPosts
}

func GetUserComments(r *http.Request) []Comment {
	user_id := GetUserId(r)
	var comments []Comment
	rows, err := DB.Query("select id, content, post_id, posting_time from comment where user_id=?", user_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, post_id int
		var content, posting_time string
		rows.Scan(&id, &content, &post_id, &posting_time)
		realusername, _ := GetUserName(user_id)
		CommentLikes := GetCommentLikes(id)
		var likes, dislikes int
		for i := 0; i < len(CommentLikes); i++ {
			if CommentLikes[i].CommentMark == "true" {
				likes++
			} else {
				dislikes++
			}
		}
		title, err := GetPostTitle(post_id)
		if err != nil {
			log.Println(err)
		}
		comment := Comment{
			CommentID:   id,
			CommentText: content,
			PostID:      post_id,
			PostTitle:   title,
			AuthorID:    user_id,
			AuthorName:  realusername,
			PostingTime: posting_time,
			Likes:       likes,
			Dislikes:    dislikes,
		}
		comments = append(comments, comment)
		sort.SliceStable(comments, func(i, j int) bool { return comments[i].CommentID > comments[j].CommentID })
	}
	fmt.Println(comments)
	return comments
}

func GetPostTitle(postID int) (string, error) {
	var title string
	err := DB.QueryRow("select title from post where id=?", postID).Scan(&title)
	if err != nil {
		log.Println(err)
		return "unknown", err
	}
	return title, nil
}

func GetUserPostsByID(userID int) []Post {
	var userPosts []Post
	allPosts := GetAllPosts()
	for i := 0; i < len(allPosts); i++ {
		if allPosts[i].Author == userID {
			userPosts = append(userPosts, allPosts[i])
		}
	}
	return userPosts
}

func GetComments(post_id int) []Comment {
	var comments []Comment
	rows, err := DB.Query("select id, content, user_id, posting_time from comment where post_id=?", post_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, user_id int
		var content, posting_time string
		rows.Scan(&id, &content, &user_id, &posting_time)
		realusername, _ := GetUserName(user_id)
		CommentLikes := GetCommentLikes(id)
		var likes, dislikes int
		for i := 0; i < len(CommentLikes); i++ {
			if CommentLikes[i].CommentMark == "true" {
				likes++
			} else {
				dislikes++
			}
		}
		comment := Comment{
			CommentID:   id,
			CommentText: content,
			PostID:      post_id,
			AuthorID:    user_id,
			AuthorName:  realusername,
			PostingTime: posting_time,
			Likes:       likes,
			Dislikes:    dislikes,
		}
		comments = append(comments, comment)
		sort.SliceStable(comments, func(i, j int) bool { return comments[i].CommentID > comments[j].CommentID })
	}
	return comments
}

func GetLikes(post_id int) []UserLike {
	var UserLikes []UserLike
	rows, err := DB.Query("select id, mark, user_id from userlike where post_id=?", post_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, user_id int
		var mark string
		rows.Scan(&id, &mark, &user_id)
		realusername, _ := GetUserName(user_id)
		UserLike := UserLike{
			LikeID:     id,
			Mark:       mark,
			PostID:     post_id,
			AuthorID:   user_id,
			AuthorName: realusername,
		}
		UserLikes = append(UserLikes, UserLike)
	}
	return UserLikes
}

func GetCommentLikes(comment_id int) []CommentUserLike {
	var CommentUserLikes []CommentUserLike
	rows, err := DB.Query("select id, mark, user_id from userlike_comment where comment_id=?", comment_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, user_id int
		var mark string
		rows.Scan(&id, &mark, &user_id)
		realusername, _ := GetUserName(user_id)
		CommentUserLike := CommentUserLike{
			CommentLikeID: id,
			CommentMark:   mark,
			CommentID:     comment_id,
			UserID:        user_id,
			AuthorName:    realusername,
		}
		CommentUserLikes = append(CommentUserLikes, CommentUserLike)
	}
	return CommentUserLikes
}

func GetUserLikes(user_id int) []int {
	var likedPosts []int
	rows, err := DB.Query("select post_id from userlike where user_id=? AND mark='true'", user_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var post_id int
		rows.Scan(&post_id)
		likedPosts = append(likedPosts, post_id)
	}
	return likedPosts
}

func GetSinglePost(postID int) (int, error) {
	var id int
	err := DB.QueryRow("select id from post where id=?", postID).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func GetCategoryNameById(category_id int) (string, error) {
	var categoryName string
	err := DB.QueryRow("select name from category where id=?", category_id).Scan(&categoryName)
	if err != nil {
		log.Println(err)
	}
	return categoryName, err
}

func GetUserName(user_id int) (string, error) {
	var username string
	err := DB.QueryRow("select name from user where id=?", user_id).Scan(&username)
	if err != nil {
		fmt.Println("1before no rows", err)
		return "unknown", err
	}
	return username, nil
}

func GetUserId(r *http.Request) int {
	var user_id int
	cookie, err := r.Cookie("authenticated")
	err2 := DB.QueryRow("select user_id from session where cookievalue=?", cookie.Value).Scan(&user_id)
	if err2 != nil {
		fmt.Println("3before no rows", err, err2)
		return 0
	}
	return user_id
}

func (c CommentUserLike) CreateCommentLike() error {
	insertUserSQL := `INSERT INTO userlike_comment(mark, user_id, comment_id) VALUES (?, ?, ?)`
	statement, err := DB.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	if c.CommentID != 0 {
		_, err = statement.Exec(c.CommentMark, c.UserID, c.CommentID)
	}
	statement.Close()

	if err != nil {
		return err
	}
	return nil
}

func SubmitReport(post_id, user_id int, status string) error {
	insertUserSQL := `INSERT INTO reports(post_id, user_id, status)
	SELECT * FROM (SELECT $1 AS post_id, $2 AS user_id, $3 AS status) AS temp
	WHERE NOT EXISTS (
		SELECT post_id FROM reports WHERE post_id = $1
	) LIMIT 1;`
	statement, err := DB.Prepare(insertUserSQL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = statement.Exec(post_id, user_id, status)
	statement.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SubmitReplyReport(post_id int, admin_reply string) error {
	insertUserSQL := `UPDATE reports SET admin_reply=$1 WHERE post_id=$2;`
	statement, err := DB.Prepare(insertUserSQL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = statement.Exec(admin_reply, post_id)
	statement.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func RequestPromotion(user_id int) error {
	insertUserSQL := `INSERT INTO requests(user_id)
	SELECT * FROM (SELECT $1 AS user_id) AS temp
	WHERE NOT EXISTS (
		SELECT user_id FROM requests WHERE user_id = $1
	) LIMIT 1;`
	statement, err := DB.Prepare(insertUserSQL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = statement.Exec(user_id)
	statement.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
