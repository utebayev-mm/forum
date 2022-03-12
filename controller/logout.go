package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"net/http"
)

// Logout : used to delete user's cookies
func Logout(w http.ResponseWriter, r *http.Request) {
	// userData := Data{
	// 	Authenticated: false,
	// 	Posts:         model.GetAllPosts(),
	// }
	ctxData := r.Context().Value(utils.ContextUserKey)
	userData := ctxData.(utils.ContextUserData)
	userData.Posts = model.GetAllPosts()
	if utils.IsAuthenticated(r) {
		userData.Authenticated = true
	}
	if userData.Authenticated {
		c := http.Cookie{
			Name:   "authenticated",
			MaxAge: -1}
		http.SetCookie(w, &c)
		userID := utils.GetUserId(r)
		if utils.IsUserHaveSession(userID) {
			err := utils.DeleteSession(userID)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(userID, "cookie deleted")
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
