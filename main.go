package main

import (
	"database/sql"
	"fmt"
	"forum/controller"
	"forum/model"
	"forum/utils"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

var routes = []route{
	newRoute("GET", "/", utils.Middleware(controller.MainPage, false)),
	newRoute("POST", "/", utils.Middleware(controller.MainPage, false)),

	newRoute("POST", "/googleauth", utils.Middleware(controller.GoogleAuth, false)),
	newRoute("GET", "/googleauth", utils.Middleware(controller.GoogleAuth, false)),

	newRoute("POST", "/githubauth", utils.Middleware(controller.GithubAuth, false)),
	newRoute("GET", "/githubauth", utils.Middleware(controller.GithubAuth, false)),

	newRoute("GET", "/static/?(.*)",
		http.StripPrefix("/static", http.FileServer(http.Dir("./static"))).ServeHTTP),

	newRoute("GET", "/sign_up?(.*)", utils.Middleware(controller.SignUp, true)),
	newRoute("POST", "/sign_up?(.*)", utils.Middleware(controller.SignUp, true)),

	newRoute("GET", "/login?(.*)", utils.Middleware(controller.Login, true)),
	newRoute("POST", "/login?(.*)", utils.Middleware(controller.Login, true)),

	newRoute("GET", "/createnewpost", utils.Middleware(controller.CreateNewPost, false)),
	newRoute("POST", "/createnewpost", utils.Middleware(controller.CreateNewPost, false)),

	newRoute("GET", "/myposts", utils.Middleware(controller.MyPosts, false)),
	newRoute("POST", "/myposts", utils.Middleware(controller.MyPosts, false)),

	newRoute("GET", "/manageusers", utils.Middleware(controller.ManageUsers, false)),
	newRoute("POST", "/manageusers", utils.Middleware(controller.ManageUsers, false)),

	newRoute("GET", "/managecategories", utils.Middleware(controller.ManageCategories, false)),
	newRoute("POST", "/managecategories", utils.Middleware(controller.ManageCategories, false)),

	newRoute("GET", "/reports", utils.Middleware(controller.ManageReports, false)),
	newRoute("POST", "/reports", utils.Middleware(controller.ManageReports, false)),

	newRoute("GET", "/createacategory", utils.Middleware(controller.CreateACategory, false)),
	newRoute("POST", "/createacategory", utils.Middleware(controller.CreateACategory, false)),

	newRoute("GET", "/mycomments", utils.Middleware(controller.MyComments, false)),
	newRoute("POST", "/mycomments", utils.Middleware(controller.MyComments, false)),

	newRoute("GET", "/myposts/edit/?(.*)", utils.Middleware(controller.EditMyPost, false)),
	newRoute("POST", "/myposts/edit/?(.*)", utils.Middleware(controller.EditMyPost, false)),

	newRoute("GET", "/mycomments/edit/?(.*)", utils.Middleware(controller.EditMyComment, false)),
	newRoute("POST", "/mycomments/edit/?(.*)", utils.Middleware(controller.EditMyComment, false)),

	newRoute("GET", "/user/?(.*)", utils.Middleware(controller.UserPosts, false)),

	newRoute("GET", "/logout", utils.Middleware(controller.Logout, false)),

	newRoute("GET", "/post/([0-9]+)", utils.Middleware(controller.Post, false)),
	newRoute("POST", "/post/([0-9]+)", utils.Middleware(controller.Post, false)),

	newRoute("GET", "/category/([0-9]+)", utils.Middleware(controller.Category, false)),

	newRoute("GET", "/tags/?(.*)", utils.Middleware(controller.Tag, false)),

	newRoute("GET", "/likedposts?(.*)", utils.Middleware(controller.LikedPosts, false)),

	newRoute("POST", "/like/", utils.Middleware(controller.Like, false)),
	newRoute("GET", "/view_all_notification/", utils.Middleware(controller.ViewAllNotification, false)),
	// newRoute("GET", "/profile/", utils.Middleware(controller.ProfilePage, false)),
	newRoute("GET", "/profile/", utils.Middleware(controller.ProfilePage, false)),
	newRoute("POST", "/profile/", utils.Middleware(controller.ProfilePage, false)),
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string

	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			route.handler(w, r)
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		controller.ErrorHandler(r, w, 405)
		return
	}
	controller.ErrorHandler(r, w, 404)

}

func main() {
	db, err := sql.Open("sqlite3", "file:database.db?sqlite_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	db.Ping()
	defer db.Close()

	initDb(db)
	model.DB = db

	Addr := ":8080"

	fmt.Println("server started at port", Addr)

	router := http.HandlerFunc(Serve)
	http.ListenAndServe(Addr, router)
}

// initDb - create all tables
func initDb(db *sql.DB) {
	userQuery, err := os.ReadFile("model/sql/initial.sql")
	if err != nil {
		log.Println(err)
	}
	if _, err := db.Exec(string(userQuery)); err != nil {
		log.Println(err)
	}
}
