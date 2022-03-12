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

// MyPosts : to show the posts the user published
func MyPosts(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/myposts.html",
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
			postToDelete := r.FormValue("delete")
			postToDeleteID, _ := strconv.Atoi(postToDelete)
			model.DeleteCommentInMyComments(postToDeleteID)
			model.DeletePostInActivity(postToDeleteID)
			fmt.Println("post to delete id=", postToDelete)
			err := model.DeletePost(postToDeleteID)
			if err != nil {
				log.Println(err)
			}
			err = model.DeleteLikeAfterPost(postToDeleteID)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(r.URL.Path)
			http.Redirect(w, r, "/myposts", http.StatusSeeOther)
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
		userData.Posts = model.GetUserPosts(r)
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	} else {
		ErrorHandler(r, w, http.StatusUnauthorized)
	}
}
