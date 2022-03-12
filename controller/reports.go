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

type ReportsData struct {
	Authenticated bool
	Reports       []model.Report
	Username      string
	CurrentUrl    string
	Notifications []utils.UserNotification
	UserRole      int
}

func ManageReports(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/reports.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	userData := ReportsData{
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
			userData.Reports = model.GetAllReports()
			if r.Method == http.MethodPost {
				if r.FormValue("replytext") != "" && r.FormValue("reportreply") != "" {
					reportreplyID, err := strconv.Atoi(r.FormValue("reportreply"))
					if err != nil {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					replyText := r.FormValue("replytext")
					errReplyText := CheckSpaces(replyText)
					if errReplyText == "empty" {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					model.SubmitReplyReport(reportreplyID, replyText)
				} else if r.FormValue("changestatus") != "" {
					postChecked, err := strconv.Atoi(r.FormValue("changestatus"))
					if err != nil {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					errReport := model.DeleteReport(postChecked)
					if errReport != nil {
						fmt.Println(errReport)
					}
					fmt.Println("post checked: ", postChecked)
				}
				http.Redirect(w, r, "/reports", http.StatusSeeOther)
			}
		} else if userData.UserRole == 2 {
			username, errName := model.GetUserName(userID)
			if errName != nil {
				fmt.Println(errName)
				return
			}
			userData.Username = username
			userData.Notifications = utils.GetAllUnviewUserNotification(model.GetUserId(r))
			userData.Reports = model.GetModeratorReports(userID)
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
