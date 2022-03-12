package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/model"
	"forum/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const clientID = "9c0298d703179d64b5c8"
const clientSecret = "e8e676dfa27235b3f6c2007df4a91a23cecfc0a6"

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

type GitHubStruct struct {
	Login string `json:"login"`
}

func GithubAuth(w http.ResponseWriter, r *http.Request) {
	// First, we need to get the value of the `code` query param
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := r.FormValue("code")

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// bytes, err := ioutil.ReadAll(req.Body)
	// fmt.Println(bytes)

	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	var httpClient http.Client
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err, "\n")
		w.WriteHeader(http.StatusBadRequest)
	}

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token

	fmt.Println("github access token ", t.AccessToken)

	// curl -H "Authorization: token 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" https://api.github.com/user

	req2, err2 := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err2 != nil {
		// handle err
	}
	req2.Header.Set("Authorization", "token "+t.AccessToken)

	resp, err3 := http.DefaultClient.Do(req2)
	if err3 != nil {
		// handle err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	GitHubStruct := &GitHubStruct{}

	err4 := json.Unmarshal(bytes, &GitHubStruct)
	if err4 != nil {
		log.Println(err4)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println(GitHubStruct.Login)

	Newuser := model.User{
		Name:    GitHubStruct.Login,
		Email:   GitHubStruct.Login + "@github",
		Role_id: 1,
	}
	err = Newuser.SetUserIdByEmail()
	if err == sql.ErrNoRows {
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
	} else {
		fmt.Println("google user already exists")
		Newuser.Id = Newuser.UserIdByEmail()
		if Newuser.IsUserNeedCookie() {
			utils.AddSession(w, Newuser.Id)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
