package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/johnietre/my-site/server/apps"
	"github.com/johnietre/my-site/server/blogs"
	"github.com/johnietre/my-site/server/repos"
	utils "github.com/johnietre/utils/go"
)

var (
	tmplsDir, remoteIP string
	baseTmplPath       string

	adminUsername, adminPassword string

	tmplNames = []string{
		"home", "me", "blog", "journal", "apps",
		"admin/login", "admin/base",
		"admin/home",
		"admin/me",
		"admin/blog",
		"admin/journal",
		"admin/apps",
		"admin/apps-issues", "admin/apps-issues-reply",
		"admin/apps-list", "admin/apps-list-edit",
		"admin/site",
	}
	tmpls = utils.NewSyncMap[string, *template.Template]()
)

func InitHandlers(tmplsDirPath, remIP string, aConfig AdminConfig) error {
	tmplsDir, remoteIP = tmplsDirPath, remIP
	adminConfig = aConfig
	baseTmplPath = filepath.Join(tmplsDir, "base.tmpl")

	for _, name := range tmplNames {
		if err := loadTmpl(name); err != nil {
			return err
		}
	}
	return nil
}

func CreateRouter(staticDir string) http.Handler {
	// Populate the routes
	router := http.NewServeMux()

	static := http.FileServer(http.Dir(staticDir))
	router.Handle("/static/", http.StripPrefix("/static", static))

	for _, name := range tmplNames {
		if err := loadTmpl(name); err != nil {
			log.Fatal(err)
		}
	}
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/me", meHandler)
	router.HandleFunc("/blog", blogHandler)
	router.HandleFunc("/journal", journalHandler)
	router.HandleFunc("/apps", appsHandler)

	router.Handle(
		"/admin/",
		http.StripPrefix(
			"/admin",
			http.Handler(adminAuthMiddleware(http.HandlerFunc(adminHandler))),
		),
	)

	return router
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.Redirect(w, r, "", http.StatusFound)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "home", Data: repos.NewReposPageData()}
	execTmpl("home", w, data)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "me"}
	execTmpl("me", w, data)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if id := query.Get("id"); id != "" {
		return
	}
	data := PageData{Active: "blog", Data: blogs.NewBlogsPageData()}
	execTmpl("blog", w, data)
}

func appsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		data := PageData{Active: "apps", Data: apps.NewAppsPageData()}
		execTmpl("apps", w, data)
		return
	}
	if err := r.ParseForm(); err != nil {
		// TODO: Error and response codes
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	issue := apps.AppIssue{
		Email:       r.PostFormValue("email"),
		Reason:      r.PostFormValue("reason"),
		Subject:     r.PostFormValue("subject"),
		Description: r.PostFormValue("description"),
		Ip:          ip,
		Timestamp:   time.Now().Unix(),
	}
	_, err := apps.AddAppIssue(r.PostFormValue("app"), issue)
	if err != nil {
		msg, status := "", 0
		if apps.IsAAIError(err) {
			msg, status = err.Error(), http.StatusBadRequest
		} else {
			log.Printf("error adding app issue: %v", err)
			msg, status = "Internal server error", http.StatusInternalServerError
		}
		//status = 200
		http.Error(
			w,
			fmt.Sprintf(
				`<p class="apps-error-resp" style="color:red">Error: %s</p>`,
				msg,
			),
			status,
		)
		return
	}
	w.Write([]byte(`<p>Success</p>`))
}

func journalHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "journal"}
	execTmpl("journal", w, data)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Login
	path := r.URL.Path
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	// Parts will never have len == 0 since `sep` (second arg) isn't empty
	switch parts[0] {
	case "":
		/*
		   if tmpl, loaded := tmpls.Load("admin/base"); !loaded {
		     log.Printf("missing admin/base template")
		     http.Error(w, "Internal server error", http.StatusInternalServerError)
		   } else if err := tmpl.Execute(w, PageData{}); err != nil {
		     log.Printf("error executing admin/base template: %v", err)
		   }
		*/
		execTmpl("admin/base", w, PageData{})
		return
	case "home":
		handleAdminHome(w, r, parts[1:])
	case "me":
		handleAdminMe(w, r, parts[1:])
	case "blog":
		handleAdminBlog(w, r, parts[1:])
	case "journal":
		handleAdminJournal(w, r, parts[1:])
	case "apps":
		handleAdminApps(w, r, parts[1:])
	case "site":
		handleAdminSite(w, r, parts[1:])
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
	/*
		if path == "/parse" {
			if parseTemplate(w, r) {
				w.Write([]byte("Success"))
				return
			}
		}
	*/
}

func handleAdminHome(w http.ResponseWriter, r *http.Request, parts []string) {
	execTmpl("admin/home", w, PageData{})
}

func handleAdminMe(w http.ResponseWriter, r *http.Request, parts []string) {
	execTmpl("admin/me", w, PageData{})
}

func handleAdminBlog(w http.ResponseWriter, r *http.Request, parts []string) {
	execTmpl("admin/blog", w, PageData{})
}

func handleAdminJournal(w http.ResponseWriter, r *http.Request, parts []string) {
	execTmpl("admin/journal", w, PageData{})
}

func handleAdminApps(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 {
		execTmpl("admin/apps", w, PageData{})
		return
	}
	switch parts[0] {
	case "issues":
		handleAdminAppsIssues(w, r, parts[1:])
		return
	case "list":
		handleAdminAppsList(w, r, parts[1:])
		return
	}
}

func handleAdminAppsIssues(w http.ResponseWriter, r *http.Request, parts []string) {
	// TODO
	/*
	  partsLen := len(parts)
	  if partsLen == 0 {
	    query := r.URL.Query()
	    giq := apps.GetIssuesQuery{}
	    if val := query.Get("app"); val != "" {
	    }
	    if val := query.Get("sort-desc"); val != "" {
	    }
	    apps.GetIssues(giq)
	    return
	  }
	  r.ParseForm()
	  switch val := r.FormValue("sort-by"); val {
	  case "app":
	  case "reason":
	  case "timestamp":
	  case "":
	  default:
	    // TODO
	  }
	  r.FormValue("sort-desc")
	  r.FormValue("filter-app")
	  switch val := r.FormValue("filter-reason"); val {
	    // TODO
	  }
	  if val := r.FormValue("filter-replied-to"); val != "" {
	  }
	*/
}

func handleAdminAppsList(w http.ResponseWriter, r *http.Request, parts []string) {
	partsLen := len(parts)
	if partsLen == 0 {
		data, err := apps.NewAALPageData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error getting admin/apps/list page data: %v", err)
		} else {
			execTmpl("admin/apps-list", w, PageData{Data: data})
		}
		return
	}
	if partsLen != 1 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	switch parts[0] {
	case "0":
		handleAdminAppsListNew(w, r)
		return
	case "reload":
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		} else if err := apps.LoadAppData(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write([]byte("Success"))
		}
		return
	}
	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	handleAdminAppsListEdit(w, r, id)
}

func handleAdminAppsListEdit(w http.ResponseWriter, r *http.Request, id uint64) {
	// TODO
	data, err := apps.NewAALEPageData(id)
	if r.Method == http.MethodPut {
		//
		return
	} else {
		if err != nil {
			if err == apps.ErrNotFound {
				http.Error(
					w,
					fmt.Sprintf("app with ID %d not found", id),
					http.StatusNotFound,
				)
			} else {
				log.Printf("error getting app data: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	execTmpl("admin/apps-list-edit", w, PageData{Data: data})
}

func handleAdminAppsListNew(w http.ResponseWriter, r *http.Request) {
	data, _ := apps.NewAALEPageData(0)
	if r.Method == http.MethodPost {
		onAppStore, err := strconv.ParseBool(r.PostFormValue("on-app-store"))
		if err != nil {
			http.Error(w, "invalid on-app-store value", http.StatusBadRequest)
			return
		}
		onPlayStore, err := strconv.ParseBool(r.PostFormValue("on-play-store"))
		if err != nil {
			http.Error(w, "invalid on-play-store value", http.StatusBadRequest)
			return
		}
		app := apps.App{
			Name:        r.PostFormValue("name"),
			Description: r.PostFormValue("description"),
			Webpage:     r.PostFormValue("webpage"),
			OnAppStore:  onAppStore,
			OnPlayStore: onPlayStore,
		}
		if app.Name == "" {
			http.Error(w, "missing app name", http.StatusBadRequest)
			return
		}
		app, err = apps.AddApp(app)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.App = app
	}
	execTmpl("admin/apps-list-edit", w, PageData{Data: data})
}

type AdminSiteData struct {
	TmplNames []string
}

func handleAdminSite(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 {
		data := PageData{
			Data: AdminSiteData{
				TmplNames: tmplNames,
			},
		}
		execTmpl("admin/site", w, data)
		return
	}
	if parts[0] == "parse" {
		handleAdminSiteParse(w, r)
		return
	}
}

func handleAdminSiteParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	} else if err := r.ParseForm(); err != nil {
		// TODO
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	names := r.PostForm["names"]
	if len(names) == 0 {
		names = tmplNames
	}
	resStr := ""
	for _, name := range names {
		if err := loadTmpl(name); err != nil {
			log.Printf("error parsing %s template: %v", name, err)
			resStr += err.Error() + "\n"
		}
	}
	if resStr == "" {
		resStr = "Success!"
	}
	w.Write([]byte(resStr))
}

const (
	adminJwtIssuer = "https://johnietre.com/admin"
)

type AdminConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Cookie   struct {
		Name   string `json:"name"`
		Path   string `json:"path,omitempty"`
		Domain string `json:"domain,omitempty"`
		MaxAge int64  `json:"maxAge"`
		Secure bool   `json:"secure,omitempty"`
	} `json:"cookie"`
	Jwt struct {
		Issuer     string `json:"issuer"`
		Timeout    int64  `json:"timeout"`
		SigningKey string `json:"signingKey"`
	} `json:"jwt"`
}

var (
	adminConfig AdminConfig
)

type AdminJWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "" && path[0] == '/' {
			path = path[1:]
		}
		if path == "login" {
			adminLoginHandler(w, r)
			return
		}
		cookie, _ := r.Cookie(adminConfig.Cookie.Name)
		if path != "logout" {
			if validateAdminJwtCookie(cookie) {
				next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}
		cookie = newAdminJwtCookie("")
		cookie.Expires, cookie.MaxAge = time.Time{}, -1
		http.SetCookie(w, cookie)
		if path == "" || path == "/" {
			execTmpl("admin/login", w, PageData{})
			return
		} else if path == "logout" {
			w.Header().Set("Location", "../admin")
			w.WriteHeader(http.StatusFound)
			return
		}
		w.Write([]byte("Unauthorized"))
	})
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if username != adminConfig.Username || password != adminConfig.Password {
		w.WriteHeader(http.StatusUnauthorized)
		execTmpl("admin/login", w, PageData{})
		return
	}
	tokStr, err := newAdminJwt(username)
	if err != nil {
		log.Printf("eror creating admin jwt: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, newAdminJwtCookie(tokStr))
	//execTmpl("admin/base", w, PageData{})
	w.Header().Set("Location", "../admin")
	w.WriteHeader(http.StatusFound)
}

func newAdminJwtCookie(value string) *http.Cookie {
	now := time.Now()
	return &http.Cookie{
		Name:     adminConfig.Cookie.Name,
		Value:    value,
		Path:     adminConfig.Cookie.Path,
		Domain:   adminConfig.Cookie.Domain,
		Expires:  now.Add(time.Duration(adminConfig.Cookie.MaxAge) * time.Second),
		Secure:   adminConfig.Cookie.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func newAdminJwt(username string) (string, error) {
	now := time.Now()
	claims := &AdminJWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: adminConfig.Jwt.Issuer,
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(adminConfig.Jwt.Timeout) * time.Second),
			),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(adminConfig.Jwt.SigningKey))
}

func validateAdminJwtCookie(cookie *http.Cookie) bool {
	return cookie != nil && validateAdminJwt(cookie.Value)
}

func validateAdminJwt(tokStr string) bool {
	tok, err := jwt.ParseWithClaims(
		tokStr,
		&AdminJWTClaims{},
		func(tok *jwt.Token) (any, error) {
			if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
			}
			return []byte(adminConfig.Jwt.SigningKey), nil
		},
	)
	if err != nil {
		return false
	}
	claims, ok := tok.Claims.(*AdminJWTClaims)
	if !ok {
		return false
	}
	return tok.Valid &&
		claims.Username == adminConfig.Username &&
		claims.RegisteredClaims.Issuer == adminConfig.Jwt.Issuer
}

// Returns true if successful
func parseTemplate(w http.ResponseWriter, r *http.Request) bool {
	// TODO: Return Error?
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || host != remoteIP {
		// TODO
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<p style="color:red">Unauthorized</p>`))
		return false
	}
	for _, name := range tmplNames {
		if err := loadTmpl(name); err != nil {
			log.Printf("error parsing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return false
		}
	}
	log.Printf("parsed template")
	return true
}

func execTmpl(name string, w http.ResponseWriter, data PageData) {
	if tmpl, loaded := tmpls.Load(name); !loaded {
		log.Printf("missing %s template", name)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else if err := tmpl.Execute(w, data); err != nil {
		log.Printf("error executing %s template: %v", name, err)
		// NOTE: This could result in probable double-write
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func loadTmpl(tmplName string) error {
	var tmpl *template.Template
	var err error
	if strings.HasPrefix(tmplName, "admin/") {
		tmpl, err = template.ParseFiles(filepath.Join(tmplsDir, tmplName+".tmpl"))
	} else {
		tmpl, err = template.ParseFiles(
			baseTmplPath, filepath.Join(tmplsDir, tmplName+".tmpl"),
		)
	}
	if err != nil {
		return fmt.Errorf("error parsing %s tmpl file: %v", tmplName, err)
	}
	// Check to make sure the template executes without err
	if err := tmpl.Execute(io.Discard, PageData{}); err != nil {
		return fmt.Errorf("error executing %s template: %v", tmplName, err)
	}
	tmpls.Store(tmplName, tmpl)
	return nil
}

type PageData struct {
	Active string
	Data   any
}
