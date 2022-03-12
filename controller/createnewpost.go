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

// CreateNewPost : sends the INSERT query into the post table of the database
func CreateNewPost(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./static/templates/createnewpost.html",
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
	userData.Notifications = utils.GetAllUnviewUserNotification(userData.UserId)

	if r.Method == http.MethodPost {
		if userData.Authenticated {

			PostTitle := r.FormValue("Posttitle")
			PostCategory := r.FormValue("Postcategory")
			PostContent := r.FormValue("Postcontent")
			categoryID := utils.GetCategoryId(PostCategory)
			userID := utils.GetUserId(r)
			PostTags := r.FormValue("Posttags")
			if !utils.SpecialCharacterValidator(PostTags) {
				ErrorHandler(r, w, http.StatusBadRequest)
				return
			}
			fmt.Println("posttags", PostTags)
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
				fmt.Printf("Successfully Uploaded File\n")
				newPost.Image = ("static/media/" + finalFileName)
			} else {
				newPost.Image = ("/")
			}
			err := newPost.Create()
			if err != nil {
				log.Println(err)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
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

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func CheckSpaces(PostTitle string) string {
	counter := 0
	for _, letter := range PostTitle {
		if letter == 32 || letter == 10 || letter == 13 {
			counter++
		}
	}
	if len(PostTitle) == counter {
		return "empty"
	}
	return "not empty"
}
