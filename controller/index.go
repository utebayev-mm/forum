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

// MainPage : shows all posts
func MainPage(w http.ResponseWriter, r *http.Request) {
	// TODO: should handle route and method errors
	templates := []string{
		"./static/templates/index.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)

	if err != nil {
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}

	ctxData := r.Context().Value(utils.ContextUserKey)
	userData := ctxData.(utils.ContextUserData)
	userData.Posts = model.GetAllPosts()
	userData.Notifications = utils.GetAllUnviewUserNotification(userData.UserId)

	if r.Method == http.MethodPost {
		if userData.Authenticated == true {
			postToDelete := r.FormValue("delete")
			postToDeleteID, _ := strconv.Atoi(postToDelete)
			fmt.Println("post to delete id=", postToDelete)
			model.DeleteCommentInMyComments(postToDeleteID)
			model.DeletePostInActivity(postToDeleteID)
			model.DeletePostInReports(postToDeleteID)
			err := model.DeletePost(postToDeleteID)
			if err != nil {
				log.Println(err)
			}
			err = model.DeleteLikeAfterPost(postToDeleteID)
			if err != nil {
				log.Println(err)
			}
			err = model.DeletePostInActivity(postToDeleteID)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(r.URL.Path)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			ErrorHandler(r, w, http.StatusUnauthorized)
		}
	}
	if err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
	if err := tmpl.Execute(w, userData); err != nil {
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}
}
