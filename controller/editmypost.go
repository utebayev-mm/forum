package controller

import (
	"fmt"
	"forum/model"
	"forum/utils"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// EditPostData : stores the data to execute
// type EditPostData struct {
// 	Authenticated bool
// 	Post          model.Post
// 	Username      string
// }

// EditMyPost : sends the UPDATE query into the database with the renewed post information
func EditMyPost(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/editpost.html",
		"./static/templates/navbar.html",
		"./static/templates/footer.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err)
		ErrorHandler(r, w, http.StatusInternalServerError)
		return
	}
	strID := r.URL.Path[14:]
	id, _ := strconv.Atoi(strID)
	fmt.Println("post id to edit:", id)
	userData := PostData{
		Authenticated: false,
		Post:          model.GetPostToEdit(id),
	}
	userID := utils.GetUserId(r)
	username, errName := model.GetUserName(userID)
	if errName != nil {
		fmt.Println(errName)
		return
	}
	userData.Username = username
	fmt.Println("post", userData.Post)

	if utils.IsAuthenticated(r) {
		userData.Authenticated = true
		userData.UserRole = utils.GetUserRole(userID)
	}
	if r.Method == http.MethodPost {
		if userData.Authenticated {
			PostTitle := r.FormValue("Posttitle")
			PostCategory := r.FormValue("Postcategory")
			PostContent := r.FormValue("Postcontent")
			categoryID := utils.GetCategoryId(PostCategory)
			userID := utils.GetUserId(r)
			PostTags := r.FormValue("Posttags")
			PostImage := userData.Post.Image
			if !utils.SpecialCharacterValidator(PostTags) {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			dt := time.Now()
			errSpacesTitle := CheckSpaces(PostTitle)
			errSpacesCategory := CheckSpaces(PostCategory)
			errSpacesContent := CheckSpaces(PostContent)
			errSpacesTags := CheckSpaces(PostTags)
			if errSpacesTitle == "empty" || errSpacesCategory == "empty" || errSpacesContent == "empty" || errSpacesTags == "empty" {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			newPost := model.Post{
				Title:       PostTitle,
				Body:        PostContent,
				Author:      userID,
				Category:    categoryID,
				Tags:        PostTags,
				Image:       PostImage,
				PostingTime: dt.Format("01-02-2006 15:04:05"),
			}
			if len(PostTitle) == 0 || len(PostContent) == 0 {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}

			file, handler, errFile := r.FormFile("postImage")
			if file != nil && handler != nil {

				if err != nil {
					log.Println(err)
					return
				}
				if errFile != nil {
					fmt.Println("Error Retrieving the File")
					fmt.Println(err)
					return
				}
				defer file.Close()
				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				if handler.Size > 20971520 {
					ErrorHandler(r, w, http.StatusRequestEntityTooLarge)
					return
				}
				fmt.Printf("MIME Header: %+v\n", handler.Header)
				randomNumber := strconv.Itoa(rand.Intn(100000))
				finalFileName := randomNumber + handler.Filename
				dst, errFileCreate := os.Create("static/media/" + finalFileName)
				defer dst.Close()
				if errFileCreate != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if _, errFileCopy := io.Copy(dst, file); err != nil {
					http.Error(w, errFileCopy.Error(), http.StatusInternalServerError)
					return
				}
				f, err := os.Open("static/media/" + finalFileName)
				if err != nil {
					ErrorHandler(r, w, http.StatusBadRequest)
					return
				}
				defer f.Close()

				// Get the content
				contentType, err := GetFileContentType(f)
				if err != nil {
					ErrorHandler(r, w, http.StatusBadRequest)
					return
				}

				fmt.Println("Content Type: " + contentType)
				if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/gif" {
					fmt.Println("invalid content type!")
					fileToDelete := os.Remove("static/media/" + finalFileName)
					if fileToDelete != nil {
						fmt.Println(fileToDelete)
					}
					ErrorHandler(r, w, http.StatusBadRequest)
					return
				}
				fmt.Println("FILE TO DELETE", userData.Post.Image)
				if len(userData.Post.Image) > 0 {
					fileToDeleteErr := os.Remove(userData.Post.Image[:])
					if fileToDeleteErr != nil {
						fmt.Println(fileToDeleteErr)
					}
				}
				fmt.Printf("Successfully Uploaded File\n")
				newPost.Image = ("static/media/" + finalFileName)
				// fmt.Println("new image", newPost.Image)
			}

			err := newPost.Update(id)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("post", id, "updated")
			http.Redirect(w, r, "/", http.StatusSeeOther)
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
		fmt.Println("error, unauthorized user")
		ErrorHandler(r, w, http.StatusUnauthorized)
		return
	}

}
