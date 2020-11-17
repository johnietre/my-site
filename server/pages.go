package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	ErrMsg string
}

var (
	pageLogger      *log.Logger
	homeTemplates   *template.Template
	stocksTemplates *template.Template
	chatTemplates   *template.Template
	convoTemplates  *template.Template
	loginTemplates *template.Template
	registerTemplates *template.Template
	tDir string
)

func init() {
	pageLogger = log.New(os.Stdout, "Page Server: ", log.LstdFlags)
	// Conditional possibly unneeded
	if strings.Contains(cwd, "server") {
		tDir = "../templates/"
	} else {
		tDir = "./templates/"
	}
	// Load the templates
	parse()
}

func parse() {
	homeTemplates = template.Must(template.ParseFiles(
		tDir+"base.html",
		tDir+"home.html"))
	stocksTemplates = template.Must(template.ParseFiles(
		tDir+"base.html",
		tDir+"stocks.html"))
	chatTemplates = template.Must(template.ParseFiles(
		tDir+"base.html",
		tDir+"chat.html"))
	convoTemplates = template.Must(template.ParseFiles(
		tDir+"base.html",
		tDir+"convo.html"))
	loginTemplates = template.Must(template.ParseFiles(
		tDir+"base.html",
		tDir+"login.html"))
	registerTemplates = template.Must(template.ParseFiles(
		tDir+"base.html",
		tDir+"register.html"))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	parse()
	if err := homeTemplates.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		pageLogger.Println(err)
	}
}

func stocksPageHandler(w http.ResponseWriter, r *http.Request) {
	parse()
	if err := stocksTemplates.ExecuteTemplate(w, "stocks.html", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		pageLogger.Println(err)
	}
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
	parse()
	if strings.HasSuffix(r.URL.Path, "/chat/") {
		if err := chatTemplates.ExecuteTemplate(w, "chat.html", nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			pageLogger.Println(err)
		}
	} else {
		if err := convoTemplates.ExecuteTemplate(w, "convo.html", nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			pageLogger.Println(err)
		}
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	parse()
	pageData := newPageData()
	loggedIn := false
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		println(email, password)
		if email != "johnietre@gmail.com" {
			pageData.ErrMsg = "Invalid email or password"
		} else {
			loggedIn = true
		}
	}
	if loggedIn {
		homePageHandler(w, r)
	} else {
		if err := loginTemplates.ExecuteTemplate(w, "login.html", pageData); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			pageLogger.Println(err)
		}
	}
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	parse()
	pageData := newPageData()
	registered := false
	if r.Method == http.MethodPost {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("email")
		password := r.FormValue("password")
		println(fname, lname, email, password)
		if email == "johnietre@gmail.com" {
			pageData.ErrMsg = "Email already assigned to account"
		} else {
			registered = true
		}
	}
	if registered {
		http.Redirect(w, r, "/", http.StatusAccepted)
	} else {
		if err := registerTemplates.ExecuteTemplate(w, "register.html", pageData); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			pageLogger.Println(err)
		}
	}
}

// Delete???
func newPageData() *PageData {
	page := &PageData{""}
	return page
}