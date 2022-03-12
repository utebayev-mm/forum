package utils

import (
	"context"
	"fmt"
	"forum/model"
	"net/http"
)

type ContextUserData struct {
	Authenticated bool
	Posts         []model.Post
	UserId        int
	CurrentUrl    string
	Username      string
	Notifications []UserNotification
	Activity      []UserActivity
	Comments      []model.Comment
	UserRole      int
}

type UserNotification struct {
	Id                int
	NotificationType  string
	NotificationValue string
	Viewed            bool
	Username          string
	PostId            int
}

type ContextKey string

const ContextUserKey ContextKey = "user"

func Middleware(handler func(w http.ResponseWriter, r *http.Request), IsAuth bool) http.HandlerFunc {

	if IsAuth {
		return func(w http.ResponseWriter, r *http.Request) {
			RCookie, err := r.Cookie("authenticated")
			if err != nil {
				handler(w, r)
				return
			}

			userCookie := RCookie.String()[14:]
			if IfCookieExists(userCookie) == true {
				fmt.Println("cookie exists in db", userCookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				DeleteCookie(w)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			UserData := ContextUserData{
				Authenticated: false,
				CurrentUrl:    r.URL.Path,
			}

			if IsAuthenticated(r) {
				UserData.Authenticated = true
				UserData.UserId = GetUserId(r)
				UserData.UserRole = GetUserRole(UserData.UserId)
				username, errName := model.GetUserName(UserData.UserId)
				if errName != nil {
					fmt.Println(errName)
					return
				}
				UserData.Username = username
				ctx := context.WithValue(r.Context(), ContextUserKey, UserData)
				handler(w, r.WithContext(ctx))
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserKey, UserData)
			handler(w, r.WithContext(ctx))
			return
		}
	}
}
