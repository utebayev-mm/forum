package utils

import (
	"fmt"
	"forum/model"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func AddSession(w http.ResponseWriter, user_id int) {

	if IsUserHaveSession(user_id) {
		err := DeleteSession(user_id)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user_id, "cookie deleted")
	}
	session := model.Session{
		CookieName:  "authenticated",
		CookieValue: GenerateCookieValue(),
		UserId:      user_id,
	}
	cookie := http.Cookie{
		Name:    session.CookieName,
		Value:   session.CookieValue,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	err := session.Create()
	if err != nil {
		log.Println(err)
	}
}

func GenerateCookieValue() string {
	return uuid.NewV4().String()
}

func DeleteCookie(rw http.ResponseWriter) {
	c := &http.Cookie{
		Name:     "authenticated",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
	
		HttpOnly: true,
	}
	
	http.SetCookie(rw, c)
}