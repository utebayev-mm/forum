package controller

import (
	"forum/utils"
	"net/http"
)

func ViewAllNotification(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value(utils.ContextUserKey)
	userData := ctxData.(utils.ContextUserData)
	if userData.Authenticated {
		utils.ViewAllNotification(userData.UserId)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ErrorHandler(r, w, http.StatusUnauthorized)
	return

}
