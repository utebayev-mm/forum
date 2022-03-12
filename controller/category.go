package controller

import (
	"forum/model"
	"forum/utils"
	"html/template"
	"net/http"
	"strconv"
)

type CategoryData struct {
	Authenticated bool
	CategoryID    int
	CategoryName  string
	Posts         []model.Post
	CurrentUrl    string
	Username      string
	Notifications []utils.UserNotification
	UserRole      int
}

func Category(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/category.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)

	allPosts := model.GetAllPosts()
	strID := r.URL.Path[len(r.URL.Path)-1:]
	id, _ := strconv.Atoi(strID)
	var catPosts []model.Post
	for i := 0; i < len(allPosts); i++ {
		if allPosts[i].Category == id {
			catPosts = append(catPosts, allPosts[i])
		}
	}
	CategoryName, err2 := model.GetCategoryNameById(id)
	userData := CategoryData{
		Authenticated: false,
		CategoryID:    id,
		CategoryName:  CategoryName,
		Posts:         catPosts,
		CurrentUrl:    r.URL.Path,
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
	if err2 != nil {
		ErrorHandler(r, w, http.StatusNotFound)
		return
	} else {
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	}
}
