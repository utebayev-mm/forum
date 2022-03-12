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

// MyComments : to show the comments the user published
func MyComments(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/mycomments.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	ctxData := r.Context().Value(utils.ContextUserKey)
	userData := ctxData.(utils.ContextUserData)

	if r.Method == http.MethodPost {
		if userData.Authenticated == true {
			if r.FormValue("delete") != "" {
				commentID, err := strconv.Atoi(r.FormValue("delete"))
				fmt.Println("comment to delete: ", commentID)
				if err != nil {
					log.Println(err)
					ErrorHandler(r, w, http.StatusBadRequest)
				}
				model.DeleteComment(commentID)
				fmt.Println("comment", commentID, "deleted")
				http.Redirect(w, r, "http://localhost:8080"+r.URL.Path, http.StatusSeeOther)
			}
		} else {
			ErrorHandler(r, w, http.StatusUnauthorized)
		}
	}
	if err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
	if userData.Authenticated {
		userData.Notifications = utils.GetAllUnviewUserNotification(userData.UserId)
		userData.Comments = model.GetUserComments(r)
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	} else {
		ErrorHandler(r, w, http.StatusUnauthorized)
	}
}
