package utils

import (
	"database/sql"
	"fmt"
	"forum/model"
	"net/http"
	"sort"
	"unicode"
)

// IsEmailExist - is function that check if given email unique
func IsEmailExist(email string) bool {
	var useremail string
	err := model.DB.QueryRow("select email from user where email=?", email).Scan(&useremail)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func DoesNameExist(name string) bool {
	var username string
	err := model.DB.QueryRow("select name from user where name=?", name).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func IfCookieExists(userCookie string) bool {
	var serverCookie string
	err := model.DB.QueryRow("select cookievalue from session where cookievalue=?", userCookie).Scan(&serverCookie)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func IsUserHaveSession(user_id int) bool {
	var cookievalue string
	err := model.DB.QueryRow("select cookievalue from session where user_id=?", user_id).Scan(&cookievalue)
	if err != nil {

		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}

	return true
}

func DeleteSession(user_id int) error {
	sqlstatement := `delete from session where user_id=$1`
	_, err := model.DB.Exec(sqlstatement, user_id)
	if err != nil {
		panic(err)
	}

	return nil
}

func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("authenticated")
	if err != nil {
		return false
	}
	if IfCookieExists(cookie.Value) == true {
		return true
	}
	// fake cookie
	return false
}

// GetCategoryId return existing category id
// if category not exists , create them return its ID
func GetCategoryId(category string) int {
	var categoryId int
	err := model.DB.QueryRow("select id from category where name=?", category).Scan(&categoryId)
	if err != nil {
		fmt.Println("2before no rows", err)
		// create new category
		CreateNewCategory(category)
		return GetCategoryId(category)

	}

	return categoryId
}

func CreateNewCategory(categoryName string) {
	sqlstatement := `insert into category(name) values ($1)`
	_, err := model.DB.Exec(sqlstatement, categoryName)
	if err != nil {
		panic(err)
	}
}

func GetUserId(r *http.Request) int {
	var user_id int
	cookie, err := r.Cookie("authenticated")
	// fmt.Println("cookie.Value", cookie.Value)
	err2 := model.DB.QueryRow("select user_id from session where cookievalue=?", cookie.Value).Scan(&user_id)
	if err2 != nil {
		fmt.Println("3before no rows", err, err2)
		return 0
	}
	return user_id
}

func GetUserRole(id int) int {
	var role_id int
	err := model.DB.QueryRow("select role_id from user where id=?", id).Scan(&role_id)
	if err != nil {
		fmt.Println("3before no rows", err)
		return 0
	}
	return role_id
}

func SpecialCharacterValidator(str string) bool {
	for _, s := range str {
		if !unicode.IsLetter(s) && !unicode.IsNumber(s) && !unicode.IsSpace(s) {
			return false
		}
	}

	return true
}

func GetAllUnviewUserNotification(user_id int) []UserNotification {
	rows, err := model.DB.Query("SELECT ua.id,notification_type,notification_value,viewed, user.name from user_activity as ua LEFT JOIN post on ua.post_id = post.id LEFT JOIN user on ua.user_id = user.id WHERE post.user_id = ? and ua.viewed = 'false'", user_id)
	var notifications []UserNotification

	if err != nil {
		fmt.Println(err)
		return notifications
	}
	defer rows.Close()
	for rows.Next() {
		var notification UserNotification
		err = rows.Scan(&notification.Id, &notification.NotificationType, &notification.NotificationValue, &notification.Viewed, &notification.Username)
		notifications = append(notifications, notification)
	}
	sort.SliceStable(notifications, func(i, j int) bool { return notifications[i].Id > notifications[j].Id })
	return notifications
}

func ViewAllNotification(user_id int) {
	notifications := GetAllUnviewUserNotification(user_id)
	if len(notifications) == 0 {
		return
	}
	fmt.Println(notifications)
	for _, notification := range notifications {
		updatePostSQL := `UPDATE user_activity SET viewed='true' WHERE id=$1;`
		statement, err := model.DB.Prepare(updatePostSQL)
		if err != nil {
			fmt.Println(err)
		}
		_, err = statement.Exec(notification.Id)
		statement.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

type UserActivity struct {
	Id            int
	ActivityType  string
	ActivityValue string
	Post          model.Post
	PostingTime   string
}

func GetAllUserActivity(user_id int) []UserActivity {
	// rows, err := model.DB.Query("SELECT ua.id,notification_type,notification_value, post.id, post.title, post.body, author.id, author.name , category.id, category.name , ua.posting_time from user_activity as ua  LEFT JOIN post on ua.post_id = post.id  LEFT JOIN user on ua.user_id = user.id  LEFT JOIN user as author on author.id = post.user_id LEFT JOIN category on category.id = post.category_id WHERE ua.user_id = ?", user_id)
	rows, err := model.DB.Query("SELECT id, notification_type, notification_value, post_id, posting_time FROM user_activity WHERE user_id = ?", user_id)

	var activities []UserActivity
	if err != nil {
		return activities
	}
	defer rows.Close()
	for rows.Next() {
		var activity UserActivity
		err = rows.Scan(&activity.Id, &activity.ActivityType, &activity.ActivityValue, &activity.Post.Id, &activity.Post.PostingTime)
		if err != nil {
			fmt.Println(err)
		}
		activity.Post = model.GetPostToEdit(activity.Post.Id)
		if activity.ActivityType == "COMMENT" {
			activity.PostingTime = model.GetActivityPostingTime(activity.Post.Id)
		} else {
			activity.PostingTime = activity.Post.PostingTime
		}
		activities = append(activities, activity)
		sort.SliceStable(activities, func(i, j int) bool { return activities[i].Id > activities[j].Id })
		fmt.Println(activities)
	}
	return activities
}
