package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"net/http"
	"text/template"
)

type ErrorData struct {
	Authenticated bool
	Status        int
	Error         string
	CurrentUrl    string
	Username      string
	Notifications []utils.UserNotification
	UserRole      int
}

// ErrorHandler : handles the errors and displays the status
func ErrorHandler(r *http.Request, w http.ResponseWriter, Status int) {
	w.WriteHeader(Status)
	templates := []string{
		"./static/templates/error.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "The following error occured:", Status)
		return
	}
	userData := ErrorData{
		Authenticated: false,
		Status:        Status,
		Error:         http.StatusText(Status),
		CurrentUrl:    "CurrentUrl",
	}
	var userID int
	if utils.IsAuthenticated(r) {
		userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))
		userData.Authenticated = true
		userID = model.GetUserId(r)
		utils.GetUserRole(userID)
		userData.Username, err = model.GetUserName(userID)
	}

	if err := tmpl.Execute(w, userData); err != nil {
		http.Error(w, "The following error occured:", Status)
		return
	}
	fmt.Println(Status)
}
