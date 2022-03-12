package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"html/template"
	"net/http"
	"strconv"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/profile.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	ctxData := r.Context().Value(utils.ContextUserKey)
	userData := ctxData.(utils.ContextUserData)
	userData.Notifications = utils.GetAllUnviewUserNotification(userData.UserId)
	if userData.Authenticated {
		activities := utils.GetAllUserActivity(userData.UserId)
		userData.Activity = activities
		fmt.Println(activities)
		if r.Method == http.MethodPost {
			userIDtoPromote, err1 := strconv.Atoi(r.FormValue("requestpromotion"))
			if err1 != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			fmt.Println("userID requested promotion: ", userIDtoPromote)
			model.RequestPromotion(userData.UserId)
		}
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	} else {
		ErrorHandler(r, w, http.StatusUnauthorized)
		return
	}

}
