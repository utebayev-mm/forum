package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

// SignUp : inserts a new user into the database
func SignUp(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/signup.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		verifyPassword := r.FormValue("password2")
		if password != verifyPassword {
			if err := tmpl.Execute(w, "Password mismatch"); err != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			return
		}
		email := r.FormValue("email")
		if utils.IsEmailExist(email) {
			if err := tmpl.Execute(w, "Given email is already in use"); err != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			return
		}
		hash := model.EncryptPassword(password)
		username := r.FormValue("name")
		if utils.DoesNameExist(username) {
			if err := tmpl.Execute(w, "Given username is already in use"); err != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			return
		}
		if len(username) == 0 || len(email) == 0 || len(password) == 0 {
			if err := tmpl.Execute(w, "Invalid email format"); err != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			return
		}
		if len(username) < 3 || len(password) < 3 {
			if err := tmpl.Execute(w, "Invalid email format"); err != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			return
		}
		if EmailCheck(email) == "invalid" {
			if err := tmpl.Execute(w, "Invalid email format"); err != nil {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			return
		}
		Newuser := model.User{
			Name:     username,
			Password: hash,
			Email:    email,
			Role_id:  1,
		}
		id, err := Newuser.Create()
		if err != nil {
			ErrorHandler(r, w, http.StatusInternalServerError)
		}
		err = Newuser.SetUserIdByEmail()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(Newuser)
		utils.AddSession(w, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if err := tmpl.Execute(w, ""); err != nil {
		ErrorHandler(r, w, http.StatusBadRequest)
		return
	}
}

func EmailCheck(email string) string {
	atCounter := 0
	dotCounter := 0
	dotIndex := 0
	for index, letter := range email {
		if letter == '@' {
			atCounter++
		}
		if letter == '.' {
			dotCounter++
			dotIndex = index + 1
		}
	}
	// fmt.Println("dot index", dotIndex)
	// fmt.Println("difference", len(email)-dotIndex)
	// fmt.Println("atcounter", atCounter)
	// fmt.Println("dotcounter", dotCounter)
	if atCounter != 1 || dotCounter != 1 || len(email)-dotIndex < 2 || len(email)-dotIndex > 3 {
		return "invalid"
	}
	return "valid"
}
