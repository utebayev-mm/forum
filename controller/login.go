package controller

import (
	"database/sql"
	"forum/model"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

// Login : receives the login data and checks if the user exists
func Login(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/login.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if len(email) == 0 || len(password) == 0 {
			ErrorHandler(r, w, http.StatusBadRequest)
			return
		}
		user := model.User{
			Email:    email,
			Password: password,
		}
		err = user.SetUserIdByEmail()
		if err == sql.ErrNoRows {
			// user doesn`t exist
			w.WriteHeader(http.StatusForbidden)
			tmpl.Execute(w, "email or password incorrect")
			return
		} else {
			user.Id = user.UserIdByEmail()
		}
		if user.IsPasswordCorrect() != true {
			w.WriteHeader(http.StatusForbidden)
			tmpl.Execute(w, "email or password incorrect")
			return
		}
		//TODO may be not need
		if user.IsUserNeedCookie() {

			utils.AddSession(w, user.Id)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if err := tmpl.Execute(w, ""); err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
}
