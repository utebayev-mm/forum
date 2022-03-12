package controller

import (
	"forum/model"
	"forum/utils"
	"log"
	"net/http"
	"strconv"
)

func Like(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value(utils.ContextUserKey)
	userData := ctxData.(utils.ContextUserData)

	if r.Method == http.MethodPost {
		if userData.Authenticated {
			userlike := r.FormValue("submit")
			AuthorID := utils.GetUserId(r)
			// if AuthorID == 0 {
			// 	ErrorHandler(r, w, http.StatusBadRequest)
			// }
			var Mark string
			var id int
			if len(userlike) > 4 && userlike[0:4] == "like" {
				Mark = "true"
				id1, err1 := strconv.Atoi(userlike[4:])
				if err1 != nil {
					ErrorHandler(r, w, http.StatusBadRequest)
					return
				}
				id = id1
			} else if len(userlike) > 7 && userlike[0:7] == "dislike" {
				Mark = "false"
				id1, err1 := strconv.Atoi(userlike[7:])
				if err1 != nil {
					ErrorHandler(r, w, http.StatusBadRequest)
					return
				}
				id = id1
			} else {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
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
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			ErrorHandler(r, w, http.StatusUnauthorized)
			return
		}
	}
}
