package controller

import (
	"forum/model"
	"forum/utils"
	"html/template"
	"net/http"
)

// TagData : stores the data for template execution
type TagData struct {
	Authenticated bool
	TagName       string
	Posts         []model.Post
	CurrentUrl    string
	Username      string
	Notifications []utils.UserNotification
	UserRole      int
}

// Tag : shows the posts tagged with a certain tag
func Tag(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/tags.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	strID := r.URL.Path[6:]
	tagPosts := model.GetTagPosts(strID)
	userData := TagData{
		Authenticated: false,
		TagName:       strID,
		Posts:         tagPosts,
		CurrentUrl:    "/tag",
	}
	if utils.IsAuthenticated(r) {
		userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))
		userData.Authenticated = true
		userID := model.GetUserId(r)
		userData.UserRole = utils.GetUserRole(userID)
		userData.Username, err = model.GetUserName(userID)
	}
	if err != nil {
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}
	if len(tagPosts) == 0 {
		ErrorHandler(r, w, http.StatusNotFound)
	} else {
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	}
}
