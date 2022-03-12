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

type ManageUsersData struct {
	Authenticated bool
	Users         []model.User
	Username      string
	CurrentUrl    string
	Notifications []utils.UserNotification
	UserRole      int
}

// MyPosts : to show the posts the user published
func ManageUsers(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/manageusers.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	userData := ManageUsersData{
		Authenticated: false,
	}
	userData.CurrentUrl = r.URL.Path
	if utils.IsAuthenticated(r) {
		userData.Authenticated = true
	}
	if userData.Authenticated == true {
		userID := utils.GetUserId(r)
		userData.UserRole = utils.GetUserRole(userID)
		if userData.UserRole == 3 {
			username, errName := model.GetUserName(userID)
			if errName != nil {
				fmt.Println(errName)
				return
			}
			userData.Username = username
			userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))
			userData.Users = model.GetAllUsers()
			requestedPromotion := model.GetAllRequests()
			for i := 0; i < len(userData.Users); i++ {
				for j := 0; j < len(requestedPromotion); j++ {
					if userData.Users[i].Id == requestedPromotion[j] {
						userData.Users[i].RequestedPromotion = true
					}
				}
			}
			fmt.Println(requestedPromotion)
			if r.Method == http.MethodPost {
				roleID := r.FormValue("changeroleid")
				roleUserID := r.FormValue("changeroleuserid")
				ManageUserID, err := strconv.Atoi(roleUserID)
				ManageRoleID, err2 := strconv.Atoi(roleID)
				if err != nil || err2 != nil {
					ErrorHandler(r, w, http.StatusBadRequest)
					return
				}
				fmt.Println("user id", ManageUserID)
				fmt.Println("new userrole id", ManageRoleID)
				errRole := model.UpdateUserRole(ManageUserID, ManageRoleID)
				if ManageRoleID == 2 {
					model.DeleteRequest(ManageUserID)
				}
				if errRole != nil {
					fmt.Println(errRole)
				}
				http.Redirect(w, r, "/manageusers", http.StatusSeeOther)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)

		}
		if err := tmpl.Execute(w, userData); err != nil {
			fmt.Println(err)
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	} else {
		ErrorHandler(r, w, http.StatusUnauthorized)
		return
	}

}
