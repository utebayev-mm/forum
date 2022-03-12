package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// UserData : stores the data for template execution
type UserData struct {
	Authenticated bool
	UserID        int
	Posts         []model.Post
	CurrentUrl    string
	Username      string
	Notifications []utils.UserNotification
	UserRole      int
}

// UserPosts : show the posts published by the certain user
func UserPosts(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/user.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	strID := r.URL.Path[6:]
	id, _ := strconv.Atoi(strID)
	fmt.Println("userID", id)
	username, notFound := model.GetUserName(id)
	fmt.Println("username", username)
	userData := UserData{
		Authenticated: false,
		UserID:        id,
		UserRole:      utils.GetUserRole(id),
		Username:      username,
		Posts:         model.GetUserPostsByID(id),
		CurrentUrl:    "/user",
	}
	if utils.IsAuthenticated(r) {
		userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))

		userData.Authenticated = true
	}
	if err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
	if notFound != nil {
		ErrorHandler(r, w, http.StatusNotFound)
	} else {
		if err := tmpl.Execute(w, userData); err != nil {
			fmt.Println("err", err)
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	}
}
