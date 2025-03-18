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
	jmux "github.com/johnietre/go-jmux"
	"github.com/johnietre/my-site/server/apps"
	"github.com/johnietre/my-site/server/blogs"
	"github.com/johnietre/my-site/server/repos"
	utils "github.com/johnietre/utils/go"
)

var (
	tmplsDir, remoteIP string
	baseTmplPath       string
	navbarTmplPath     string

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
	navbarTmplPath = filepath.Join(tmplsDir, "navbar.tmpl")

	for _, name := range tmplNames {
		if err := loadTmpl(name); err != nil {
			return err
		}
	}
	return nil
}

func CreateRouter(staticDir string) http.Handler {
	// Populate the routes
	router := jmux.NewRouter()
	//router := http.NewServeMux()

	static := http.FileServer(http.Dir(staticDir))
	router.Get(
    "/static/",
    jmux.WrapH(http.StripPrefix("/static", static)),
  ).MatchAny(jmux.MethodsGet())

	for _, name := range tmplNames {
		if err := loadTmpl(name); err != nil {
			log.Fatal(err)
		}
	}

	homeRouter := createHomeRouter()
	router.All("/", homeRouter)
	router.All("/home", homeRouter)

	router.All("/me", createMeRouter())
	router.All("/blog", createBlogRouter())
	router.All("/journal", createJournalRouter())
	router.All("/apps", createAppsRouter())

	router.All("/admin/", createAdminRouter())

	return router
}

func createHomeRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.GetFunc("/", homeHandler)
	router.GetFunc("/home", homeHandler)
	return router
}

func createMeRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.GetFunc("/", meHandler)
	return jmux.WrapH(http.StripPrefix("/me", router))
}

func createBlogRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.GetFunc("/", blogHandler)
	return jmux.WrapH(http.StripPrefix("/blog", router))
}

func createJournalRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.GetFunc("/", journalHandler)
	router.GetFunc("/{journal_id}", getJournalHandler)
	return jmux.WrapH(http.StripPrefix("/journal", router))
}

func createAppsRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.GetFunc("/", appsHandler)
	router.PostFunc("/issues", appsNewIssueHandler)
	return jmux.WrapH(http.StripPrefix("/apps", router))
}

func defaultHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.Redirect(w, r, "", http.StatusFound)
		return
	}
}

func homeHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
	path := r.URL.Path
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	if path != "" && path != "home" {
		http.NotFound(w, r)
		return
	}
	data := PageData{Active: "home", Data: repos.NewReposPageData()}
	execTmpl("home", c, data)
}

func meHandler(c *jmux.Context) {
	data := PageData{Active: "me"}
	execTmpl("me", c, data)
}

func blogHandler(c *jmux.Context) {
  /*
	query := c.Query()
	if id := query.Get("id"); id != "" {
		return
	}
  */
	data := PageData{Active: "blog", Data: blogs.NewBlogsPageData()}
	execTmpl("blog", c, data)
}

func appsHandler(c *jmux.Context) {
	data := PageData{Active: "apps", Data: apps.NewAppsPageData()}
	execTmpl("apps", c, data)
}

func appsNewIssueHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
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

func journalHandler(c *jmux.Context) {
	data := PageData{Active: "journal"}
	execTmpl("journal", c, data)
}

func getJournalHandler(c *jmux.Context) {
	data := PageData{Active: "journal"}
	execTmpl("journal", c, data)
}

func createAdminRouter() jmux.Handler {
	router := jmux.NewRouter()
	/*
	  router.All(
	    "/",
			jmux.WrapH(http.StripPrefix(
				"/admin",
				jmux.ToHTTP(jmux.Handler(
	        adminAuthMiddleware(jmux.HandlerFunc(adminHandler)),
	      )),
			)),
	  )
	*/
	//return router

	router.GetFunc("/home", adminHomeHandler)
	router.GetFunc("/me", adminMeHandler)
	router.GetFunc("/blog", adminBlogHandler)
	router.GetFunc("/journal", adminJournalHandler)

	router.GetFunc("/apps", adminAppsHandler)
	router.GetFunc("/apps/issues", adminAppsIssuesHandler)
	router.GetFunc("/apps/list", adminAppsListHandler)
	router.PostFunc("/apps/list/0", adminAppsListNewHandler)
	//router.GetFunc("/apps/list/reload", adminAppsListReloadHandler)
	router.GetFunc("/apps/list/{app_id}", adminAppsListHandler)

	router.GetFunc("/site", adminSiteHandler)
	router.GetFunc("/site/parse", adminSiteParseHandler)

	return jmux.WrapH(http.StripPrefix(
		"/admin",
		jmux.ToHTTP(adminAuthMiddleware(router)),
	))
}

func adminHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
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
		execTmpl("admin/base", c, PageData{})
		return
	case "home":
		//handleAdminHome(c, parts[1:])
	case "me":
		//handleAdminMe(c, parts[1:])
	case "blog":
		//handleAdminBlog(c, parts[1:])
	case "journal":
		//handleAdminJournal(c, parts[1:])
	case "apps":
		//handleAdminApps(c, parts[1:])
	case "site":
		//handleAdminSite(c, parts[1:])
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

func adminHomeHandler(c *jmux.Context) {
	execTmpl("admin/home", c, PageData{})
}

func adminMeHandler(c *jmux.Context) {
	execTmpl("admin/me", c, PageData{})
}

func adminBlogHandler(c *jmux.Context) {
	execTmpl("admin/blog", c, PageData{})
}

func adminJournalHandler(c *jmux.Context) {
	execTmpl("admin/journal", c, PageData{})
}

func adminAppsHandler(c *jmux.Context) {
	execTmpl("admin/apps", c, PageData{})
}

func adminAppsIssuesHandler(c *jmux.Context) {
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

func adminAppsListHandler(c *jmux.Context) {
  /*
	w, r := c.Writer, c.Request
	partsLen := len(parts)
	if partsLen == 0 {
		data, err := apps.NewAALPageData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error getting admin/apps/list page data: %v", err)
		} else {
			execTmpl("admin/apps-list", c, PageData{Data: data})
		}
		return
	}
	if partsLen != 1 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	switch parts[0] {
	case "0":
		handleAdminAppsListNew(c)
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
	handleAdminAppsListEdit(c, id)
  */
}

func adminAppsListEditHandler(c *jmux.Context, id uint64) {
	w, r := c.Writer, c.Request
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
	execTmpl("admin/apps-list-edit", c, PageData{Data: data})
}

func adminAppsListNewHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
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
	execTmpl("admin/apps-list-edit", c, PageData{Data: data})
}

type AdminSiteData struct {
	TmplNames []string
}

func adminSiteHandler(c *jmux.Context) {
	data := PageData{
		Data: AdminSiteData{
			TmplNames: tmplNames,
		},
	}
	execTmpl("admin/site", c, data)
}

func adminSiteParseHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
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

func adminAuthMiddleware(next jmux.Handler) jmux.Handler {
	return jmux.HandlerFunc(func(c *jmux.Context) {
		path := c.Path()
		if path != "" && path[0] == '/' {
			path = path[1:]
		}
		if path == "login" {
			adminLoginHandler(c)
			return
		}
		cookie, _ := c.Cookie(adminConfig.Cookie.Name)
		if path != "logout" {
			if validateAdminJwtCookie(cookie) {
				next.ServeC(c)
				return
			}
			c.WriteHeader(http.StatusUnauthorized)
		}
		cookie = newAdminJwtCookie("")
		cookie.Expires, cookie.MaxAge = time.Time{}, -1
		c.SetCookie(cookie)
		if path == "" || path == "/" {
			execTmpl("admin/login", c, PageData{})
			return
		} else if path == "logout" {
			c.Writer.Header().Set("Location", "../admin")
			c.WriteHeader(http.StatusFound)
			return
		}
		c.Write([]byte("Unauthorized"))
	})
}

func adminLoginHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if username != adminConfig.Username || password != adminConfig.Password {
		w.WriteHeader(http.StatusUnauthorized)
		execTmpl("admin/login", c, PageData{})
		return
	}
	tokStr, err := newAdminJwt(username)
	if err != nil {
		log.Printf("eror creating admin jwt: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, newAdminJwtCookie(tokStr))
	//execTmpl("admin/base", c, PageData{})
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
func parseTemplate(c *jmux.Context) bool {
	w, r := c.Writer, c.Request
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

func execTmpl(name string, c *jmux.Context, data PageData) {
	if tmpl, loaded := tmpls.Load(name); !loaded {
		log.Printf("missing %s template", name)
		c.InternalServerError("Internal server error")
	} else if err := tmpl.Execute(c.Writer, data); err != nil {
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
