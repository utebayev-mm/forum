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

type ManageCategoriesData struct {
	Authenticated bool
	Categories    []model.Category
	Username      string
	CurrentUrl    string
	Notifications []utils.UserNotification
	UserRole      int
}

func ManageCategories(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/managecategories.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	userData := ManageCategoriesData{
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
			userData.Categories = model.GetAllCategories()
			if r.Method == http.MethodPost {
				if r.FormValue("deletecategory") != "" {
					categoryToDelete := r.FormValue("deletecategory")
					categoryIDtoDelete, err := strconv.Atoi(categoryToDelete)
					if err != nil {
						ErrorHandler(r, w, http.StatusBadRequest)
						return
					}
					fmt.Println("category to delete id: ", categoryIDtoDelete)
					model.FindAndDeletePostsByCategory(categoryIDtoDelete)
					errCategoryPostsDelete := model.DeletePostsFromCategory(categoryIDtoDelete)
					if errCategoryPostsDelete != nil {
						fmt.Println(errCategoryPostsDelete)
					}
					errCategoryDelete := model.DeleteCategory(categoryIDtoDelete)
					if errCategoryDelete != nil {
						fmt.Println(errCategoryDelete)
					}
					http.Redirect(w, r, "/managecategories", http.StatusSeeOther)
				}
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
