package controller

import (
	"fmt"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

func CreateACategory(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/createacategory.html",
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
			CategoryTitle := r.FormValue("createacategory")
			fmt.Println("category to be created: ", CategoryTitle)
			errSpacesCategory := CheckSpaces(CategoryTitle)
			if errSpacesCategory == "empty" {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			utils.CreateNewCategory(CategoryTitle)
			http.Redirect(w, r, "/managecategories", http.StatusSeeOther)
		} else {
			ErrorHandler(r, w, http.StatusUnauthorized)
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
		ErrorHandler(r, w, http.StatusUnauthorized)
	}
}
