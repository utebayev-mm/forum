package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// EditPostData : stores the data to execute
type EditCommentData struct {
	Authenticated bool
	Comment       model.Comment
	Username      string
	CurrentUrl    string
	Notifications []utils.UserNotification
	UserRole      int
}

// EditMyPost : sends the UPDATE query into the database with the renewed post information
func EditMyComment(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/editcomment.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}
	strID := r.URL.Path[17:]
	id, _ := strconv.Atoi(strID)
	fmt.Println("comment id to edit:", id)
	userData := EditCommentData{
		Authenticated: false,
		Comment:       model.GetCommentToEdit(id),
	}
	userID := utils.GetUserId(r)
	userData.UserRole = utils.GetUserRole(userID)
	username, errName := model.GetUserName(userID)
	if errName != nil {
		fmt.Println(errName)
		return
	}
	userData.Username = username
	userData.CurrentUrl = r.URL.Path
	fmt.Println("comment", userData.Comment)

	if utils.IsAuthenticated(r) {
		userData.Authenticated = true
	}
	if r.Method == http.MethodPost {
		if userData.Authenticated {
			userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))
			CommentContent := r.FormValue("CommentContent")
			dt := time.Now()
			errSpacesCommentText := CheckSpaces(CommentContent)
			if errSpacesCommentText == "empty" {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			newComment := model.Comment{
				CommentText: CommentContent,
				PostingTime: dt.Format("01-02-2006 15:04:05"),
			}
			if len(CommentContent) == 0 {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}

			err := newComment.Update(id)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("comment", id, "updated")
			http.Redirect(w, r, "/mycomments", http.StatusSeeOther)
		}
	}
	if err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
	if userData.Authenticated {
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Println("error, unauthorized user")
		ErrorHandler(r, w, http.StatusUnauthorized)
		return
	}

}
