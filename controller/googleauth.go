package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"forum/model"
	"forum/utils"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/idtoken"
)

type GoogleTokenStruct struct {
	Key     string `json:"key"`
	IdToken string `json:"idtoken"`
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	IdTokenStruct := &GoogleTokenStruct{}

	err2 := json.Unmarshal(bytes, &IdTokenStruct)
	if err2 != nil {
		log.Println(err2)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// fmt.Println(IdTokenStruct.IdToken)
	payload, err3 := idtoken.Validate(context.Background(), IdTokenStruct.IdToken, "233275650657-5aai8oq1qpj0so9hvn58vfqofj1a1a8g.apps.googleusercontent.com")
	if err3 != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// fmt.Print("payloadClaims: ", payload.Claims, "\n")
	var name, email string
	for key, value := range payload.Claims {
		if key == "name" {
			name = fmt.Sprintf("%v", value)
		} else if key == "email" {
			email = fmt.Sprintf("%v", value)
		}
	}
	GoogleClaims, err4 := ValidateGoogleJWT(IdTokenStruct.IdToken)
	if err4 != nil {
		fmt.Println(err4)
		ErrorHandler(r, w, http.StatusForbidden)
		return
	}
	fmt.Println("token is ok")
	fmt.Println(GoogleClaims)
	Newuser := model.User{
		Name:    name,
		Email:   email,
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

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

// ValidateGoogleJWT -
func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != "233275650657-5aai8oq1qpj0so9hvn58vfqofj1a1a8g.apps.googleusercontent.com" {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}
