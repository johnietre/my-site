package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// PageData holds data for the webpage
type PageData struct {
	ErrMsg string
}

// User holds data about the user
type User struct {
	firstname string
	lastname  string
	email     string
	password  string // Hashed password
	convos []Conversation
	friends []User
}

// UserMap holds key-value pairs for user emails and User structs
type UserMap struct {
	users map[string]*User
	sync.RWMutex
}

var (
	pageLogger        *log.Logger
	homeTemplates     *template.Template
	stocksTemplates   *template.Template
	chatTemplates     *template.Template
	convoTemplates    *template.Template
	loginTemplates    *template.Template
	registerTemplates *template.Template
	tDir              string
	users             UserMap
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

	users.users = make(map[string]*User)
	users.Register("Johnie", "Rodgers", "johnietre@gmail.com", "Rj385637")
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
		// Check to see if the sender has a convo with the recipient
		// If not, start a new one
		w.Header().Set("CONTENT-TYPE", "application/json")
		writer := json.NewEncoder(w)
		// writer.Encode(msgs)
		log.Printf("%T\n", writer)

	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	parse()
	pageData := newPageData()
	loggedIn := false
	if r.Method == http.MethodPost {
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		if users.Login(email, password) {
			loggedIn = true
			println(email, "logged in")
		} else {
			pageData.ErrMsg = "Invalid email or password"
		}
	}
	if loggedIn {
		http.Redirect(w, r, "129.119.172.61:8000/stocks", http.StatusFound)
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
		fname := r.PostFormValue("fname")
		lname := r.PostFormValue("lname")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		println(fname, lname, email, password)
		if users.Register(fname, lname, email, password) {
			registered = true
			println(fname, lname, email, "registered")
		} else {
			println(email, "Nope")
			pageData.ErrMsg = "Email already assigned to account"
		}
	}
	if registered {
		http.Redirect(w, r, "/", http.StatusFound)
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

/* UserMap */

// Register checks if a user exists and registers them
func (umap *UserMap) Register(fname, lname, email, password string) bool {
	if fname == "" || lname == "" {
		return false
	}
	umap.RLock()
	user := umap.users[email]
	umap.RUnlock()
	if user != nil {
		return false
	}
	umap.Lock()
	defer umap.Unlock()
	user = umap.users[email]
	if user != nil {
		return false
	}
	hashed, err := hashPassword(password)
	if err != nil {
		pageLogger.Println(err)
		return false
	}
	user = &User{
		firstname: fname,
		lastname: lname,
		email: email,
		password: hashed}
	umap.users[email] = user
	return true
}

// Login checks if a user exists and logs them in
func (umap *UserMap) Login(email, password string) bool {
	umap.RLock()
	defer umap.RUnlock()
	user := umap.users[email]
	if user == nil {
		return false
	}
	return checkPasswordHash(password, user.password)
}
