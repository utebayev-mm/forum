package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

// LikedPosts : shows the posts the authenticated user liked
func LikedPosts(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/likedposts.html",
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

	if utils.IsAuthenticated(r) {
		userData.Authenticated = true
		userData.Notifications = utils.GetAllUnviewUserNotification(userData.UserId)
	}
	if r.URL.Path != "/likedposts" {
		if r.URL.Path != "/likedposts/" {
			ErrorHandler(r, w, http.StatusNotFound)
			return
		}
	}
	if userData.Authenticated == true {
		var likedPosts []model.Post
		userID := model.GetUserId(r)
		userLikes := model.GetUserLikes(userID)
		fmt.Println("user", userID, "liked the following posts: ", userLikes)
		allPosts := model.GetAllPosts()
		for i := 0; i < len(allPosts); i++ {
			for j := 0; j < len(userLikes); j++ {
				if allPosts[i].Id == userLikes[j] {
					likedPosts = append(likedPosts, allPosts[i])
				}
			}
		}
		userData.Posts = likedPosts
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	} else {
		ErrorHandler(r, w, http.StatusUnauthorized)
	}

	if err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
}
