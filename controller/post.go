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

// PostData : stores the data for template execution
type PostData struct {
	Authenticated bool
	Post          model.Post
	Comments      []model.Comment
	UserID        int
	CurrentUrl    string
	Username      string
	Notifications []utils.UserNotification
	UserRole      int
}

// Post : shows a single selected post
func Post(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/post.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	allPosts := model.GetAllPosts()
	strID := r.URL.Path[6:]
	id, _ := strconv.Atoi(strID)
	fmt.Println("postID", id)
	PostID, notFound := model.GetSinglePost(id)
	var post model.Post
	for i := 0; i < len(allPosts); i++ {
		if allPosts[i].Id == PostID {
			post = allPosts[i]
			break
		}
	}

	Comments := model.GetComments(id)
	userData := PostData{
		Authenticated: false,
		Post:          post,
		Comments:      Comments,
		CurrentUrl:    r.URL.Path,
	}
	fmt.Println("postID", id)
	// if userID == 0 {
	// 	ErrorHandler(r, w, http.StatusNotFound)
	// }
	if utils.IsAuthenticated(r) {
		userData.Authenticated = true
		userData.UserID = model.GetUserId(r)
		userData.UserRole = utils.GetUserRole(userData.UserID)
		userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))
		username, errName := model.GetUserName(userData.UserID)
		if errName != nil {
			fmt.Println(errName)
			return
		}
		userData.Username = username
		if userData.UserID == 0 {
			ErrorHandler(r, w, http.StatusUnauthorized)
		}
	}
	if r.Method == http.MethodPost || r.Method == http.MethodDelete {
		if userData.Authenticated {
			userData.UserID = utils.GetUserId(r)
			if r.Method == http.MethodPost {
				CommentText := r.FormValue("aComment")
				if CommentText != "" {
					AuthorID := utils.GetUserId(r)
					dt := time.Now()
					errSpacesCommentText := CheckSpaces(CommentText)
					if errSpacesCommentText == "empty" {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					newComment := model.Comment{
						CommentText: CommentText,
						PostID:      id,
						AuthorID:    AuthorID,
						PostingTime: dt.Format("01-02-2006 15:04:05"),
					}

					err := newComment.Create()
					if err != nil {
						log.Println(err)
					}

					http.Redirect(w, r, "http://localhost:8080"+r.URL.Path, http.StatusSeeOther)
				}
				userlike := r.FormValue("submit")
				if userlike != "" {
					AuthorID := utils.GetUserId(r)
					var Mark string
					var id int
					var err1 error
					if len(userlike) > 4 && userlike[0:4] == "like" {
						Mark = "true"
						id, err1 = strconv.Atoi(userlike[4:])
						if err1 != nil {
							ErrorHandler(r, w, http.StatusBadRequest)
						}
					} else if len(userlike) > 7 && userlike[0:7] == "dislike" {
						Mark = "false"
						id, err1 = strconv.Atoi(userlike[7:])
						if err1 != nil {
							ErrorHandler(r, w, http.StatusBadRequest)
						}
					} else {
						ErrorHandler(r, w, http.StatusBadRequest)
					}
					newLike := model.UserLike{
						Mark:     Mark,
						PostID:   id,
						AuthorID: AuthorID,
					}
					err := newLike.Delete()
					if err != nil {
						log.Println(err)
					}
					err2 := newLike.Create()
					if err2 != nil {
						log.Println(err2)
					}
					http.Redirect(w, r, "http://localhost:8080"+r.URL.Path, http.StatusSeeOther)
				}
				if r.FormValue("commentmark") != "" {
					mark := r.FormValue("commentmark")
					AuthorID := utils.GetUserId(r)
					if len(mark) > 11 && mark[0:11] == "commentlike" {
						commentID, err3 := strconv.Atoi(mark[11:])
						if err3 != nil {
							ErrorHandler(r, w, http.StatusBadRequest)
						}
						Mark := "true"
						fmt.Println(mark[0:11], commentID)
						newCommentLike := model.CommentUserLike{
							CommentMark: Mark,
							CommentID:   commentID,
							UserID:      AuthorID,
						}
						err := newCommentLike.Delete()
						if err != nil {
							fmt.Println(err)
						}
						err2 := newCommentLike.CreateCommentLike()
						if err2 != nil {
							fmt.Println(err2)
						}
						fmt.Println("commentlike created")
					} else if len(mark) > 14 && mark[0:14] == "commentdislike" {
						commentID, err3 := strconv.Atoi(mark[14:])
						if err3 != nil {
							ErrorHandler(r, w, http.StatusBadRequest)
						}
						fmt.Println(mark[0:14], commentID)
						Mark := "false"
						newCommentLike := model.CommentUserLike{
							CommentMark: Mark,
							CommentID:   commentID,
							UserID:      AuthorID,
						}
						err := newCommentLike.Delete()
						if err != nil {
							fmt.Println(err)
						}
						err2 := newCommentLike.CreateCommentLike()
						if err2 != nil {
							fmt.Println(err2)
						}
						fmt.Println("commentdislike created")
					} else {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					http.Redirect(w, r, "http://localhost:8080"+r.URL.Path, http.StatusSeeOther)
				}
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
				if r.FormValue("report") != "" && r.FormValue("reporttext") != "" {
					reportText := r.FormValue("reporttext")
					errSpacesReportText := CheckSpaces(reportText)
					if errSpacesReportText == "empty" {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					fmt.Println((reportText))
					reportID, errReport := strconv.Atoi(r.FormValue("report"))
					if errReport != nil {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					fmt.Println("post id to report: ", reportID)
					model.SubmitReport(reportID, userData.UserID, reportText)
				}
				if CommentText == "" && userlike == "" && r.FormValue("commentmark") == "" && r.FormValue("delete") == "" {
					// ErrorHandler(r, w, http.StatusBadRequest)
					// return
					w.WriteHeader(http.StatusBadRequest)
				}
			}
		} else {
			ErrorHandler(r, w, http.StatusUnauthorized)
			return
		}
	}
	if err != nil {
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}
	if notFound != nil {
		ErrorHandler(r, w, http.StatusNotFound)
	} else {
		if err := tmpl.Execute(w, userData); err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
			return
		}
	}

}
