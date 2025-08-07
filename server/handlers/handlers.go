package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	jmux "github.com/johnietre/go-jmux"
	"github.com/johnietre/my-site/server/blogs"
	"github.com/johnietre/my-site/server/gory-proxy"
	"github.com/johnietre/my-site/server/products"
	"github.com/johnietre/my-site/server/repos"
	"github.com/johnietre/my-site/server/sitemap"

	//goryproxy "github.com/johnietre/gory-proxy/server"
	utils "github.com/johnietre/utils/go"
)

var (
	tmplsDir, remoteIP string
	baseTmplPath       string
	navbarTmplPath     string

	adminUsername, adminPassword string

	tmplNames = []string{
		"home", "about", "blog", "journal", "products",
		"admin/login", "admin/base",
		"admin/admin",
		"admin/home",
		"admin/about",
		"admin/blog",
		"admin/journal",
		"admin/products",
		"admin/products-issues", "admin/products-issues-reply",
		"admin/products-list", "admin/products-list-edit",
		"admin/site",
	}
	tmpls = utils.NewSyncMap[string, *template.Template]()

	// TODO: replace with something that runs daily
	smUrlEntryCreators = utils.NewSyncMap[string, func(sitemap.UrlEntry) sitemap.UrlEntry]()

	proxy = goryproxy.NewRouterHandler()
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

	router.All("/about", createMeRouter())
	router.All("/blog", createBlogRouter())
	router.All("/journal", createJournalRouter())
	router.All("/products", createProductsRouter()).MatchAny(jmux.MethodsAll())

	router.All("/admin/", createAdminRouter()).MatchAny(jmux.MethodsAll())
	/* TODO: delete?
	router.GetFunc("/parse", func(c *jmux.Context) {
		if parseTemplate(c) {
			http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
			return
		}
	})
	*/

	router.GetFunc("/robots.txt", func(c *jmux.Context) {
		c.WriteFile(filepath.Join(staticDir, "robots.txt"))
	})
	router.All("/sitemap/", createSitemapRouter()).MatchAny(jmux.MethodsAll())

	router.All("/api/", createApiRouter()).MatchAny(jmux.MethodsAll())

	return router
}

func createHomeRouter() jmux.Handler {
	storeSMEntryCreator("home", 0.7)
	router := jmux.NewRouter()
	router.GetFunc("/", homeHandler)
	router.GetFunc("/home", homeHandler)
	return router
}

func createMeRouter() jmux.Handler {
	storeSMEntryCreator("about", 0.5)
	router := jmux.NewRouter()
	router.GetFunc("/", meHandler)
	return jmux.WrapH(http.StripPrefix("/about", router))
}

func createBlogRouter() jmux.Handler {
	storeSMEntryCreator("blog", 0.5)
	router := jmux.NewRouter()
	router.GetFunc("/", blogHandler)
	return jmux.WrapH(http.StripPrefix("/blog", router))
}

func createJournalRouter() jmux.Handler {
	storeSMEntryCreator("journal", 0.4)
	router := jmux.NewRouter()
	router.GetFunc("/", journalHandler)
	router.GetFunc("/{journal_id}", getJournalHandler)
	return jmux.WrapH(http.StripPrefix("/journal", router))
}

func createProductsRouter() jmux.Handler {
	storeSMEntryCreator("products", 0.8)
	router := jmux.NewRouter()
	router.GetFunc("/", productsHandler)
	router.Post(
		"/issues",
		jmux.WrapH(
			http.MaxBytesHandler(
				jmux.HandlerFunc(productsNewIssueHandler),
				1<<12,
			),
		),
	)
	return jmux.WrapH(http.StripPrefix("/products", router))
}

func createSitemapRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.AllFunc("/sitemap.xml", func(c *jmux.Context) {
		sm := sitemap.Sitemap{}
		smUrlEntryCreators.Range(func(name string, f SMEntryCreator) bool {
			sm.Urls = append(sm.Urls, f(sitemap.UrlEntry{}))
			return true
		})
		sort.Slice(sm.Urls, func(i, j int) bool {
			return sm.Urls[i].Loc < sm.Urls[j].Loc
		})
		c.RespHeader().Set("Content-Type", "application/xml")
		sitemap.NewEncoder(c.Writer).EncodeWithHeader(sm)
	})
	return jmux.WrapH(http.StripPrefix("/sitemap", router))
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
	path := cleanPath(r.URL.Path)
	if path != "" && path != "home" {
		http.NotFound(w, r)
		return
	}
	data := PageData{Active: "home", Data: repos.NewReposPageData()}
	execTmpl("home", c, data)
}

func meHandler(c *jmux.Context) {
	data := PageData{Active: "about"}
	execTmpl("about", c, data)
}

func blogHandler(c *jmux.Context) {
	/*
		query := c.Query()
		if id := query.Get("id"); id != "" {
			return
		}
		data := PageData{Active: "blog", Data: blogs.NewBlogsPageData()}
	*/
	data := PageData{Active: "blog"}
	if false {
		data = PageData{Active: "blog", Data: blogs.NewBlogsPageData()}
	}
	execTmpl("blog", c, data)
}

func productsHandler(c *jmux.Context) {
	data := PageData{Active: "products", Data: products.NewProductsPageData()}
	execTmpl("products", c, data)
}

func productsNewIssueHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
	if err := r.ParseForm(); err != nil {
		// TODO: Error and response codes
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	issue := products.ProductIssue{
		Email:       r.PostFormValue("email"),
		Reason:      r.PostFormValue("reason"),
		Subject:     r.PostFormValue("subject"),
		Description: r.PostFormValue("description"),
		Ip:          ip,
		CreatedAt:   time.Now().Unix(),
	}
	_, err := products.AddProductIssue(r.PostFormValue("product"), issue)
	if err != nil {
		msg, status := "", 0
		if products.IsAAIError(err) {
			msg, status = err.Error(), http.StatusBadRequest
		} else {
			log.Printf("error adding app issue: %v", err)
			msg, status = "Internal server error", http.StatusInternalServerError
		}
		//status = 200
		http.Error(
			w,
			fmt.Sprintf(
				`<p class="products-review-error-resp">Error: %s</p>`,
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

// TODO
func createApiRouter() jmux.Handler {
	router := jmux.NewRouter()
	router.All("/", jmux.WrapH(proxy)).MatchAny(jmux.MethodsAll())
	/*
	  router.AllFunc("/", func(c *jmux.Context) {
	    log.Print(c.Path())
	    proxy.ServeHTTP(c.Writer, c.Request)
	  }).MatchAny(jmux.MethodsAll())
	*/
	srvr := &goryproxy.Server{
		Name:   "jtgames",
		Path:   "jtgames",
		Addr:   "http://127.0.0.1:8888",
		Hidden: true,
	}
	if err := srvr.AddNewProxy("http://127.0.0.1:8888"); err != nil {
		panic(err)
	}
	if err := proxy.AddServer(srvr); err != nil {
		panic(err)
	}
	return jmux.WrapH(http.StripPrefix("/api", proxy))
	//return router
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

	router.GetFunc("/", adminHandler)
	router.GetFunc("/home", adminHomeHandler)
	router.GetFunc("/about", adminMeHandler)
	router.GetFunc("/blog", adminBlogHandler)
	router.GetFunc("/journal", adminJournalHandler)

	router.GetFunc("/products", adminProductsHandler)
	router.GetFunc("/products/list", adminProductsListHandler)
	router.PostFunc("/products/list", adminProductsNewHandler)
	//router.PostFunc("/products/list/0", adminProductsListNewHandler)
	router.GetFunc("/products/list/{product_id}", adminProductsListHandler)
	router.PutFunc("/products/list/{product_id}", adminProductsEditHandler)
	router.GetFunc("/products/issues", adminProductsIssuesHandler)
	router.GetFunc("/products/issues/{issue_id}", adminProductsIssuesHandler)
	router.PutFunc("/products/issues/{issue_id}", adminProductsIssueEditHandler)

	router.GetFunc("/site", adminSiteHandler)
	router.GetFunc("/site/parse", adminSiteParseHandler)

	router.GetFunc("/parse", func(c *jmux.Context) {
		if parseTemplate(c) {
			http.Redirect(c.Writer, c.Request, "..", http.StatusFound)
			return
		}
	})

	return jmux.WrapH(http.StripPrefix(
		"/admin",
		jmux.ToHTTP(adminAuthMiddleware(router)),
	))
}

func adminHandler(c *jmux.Context) {
	execTmpl("admin/admin", c, PageData{})
}

func _adminHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request
	// TODO: Login
	path := cleanPath(r.URL.Path)
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
	case "about":
		//handleAdminMe(c, parts[1:])
	case "blog":
		//handleAdminBlog(c, parts[1:])
	case "journal":
		//handleAdminJournal(c, parts[1:])
	case "products":
		//handleAdminProducts(c, parts[1:])
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
	execTmpl("admin/home", c, PageData{Active: "home"})
}

func adminMeHandler(c *jmux.Context) {
	execTmpl("admin/about", c, PageData{Active: "about"})
}

func adminBlogHandler(c *jmux.Context) {
	execTmpl("admin/blog", c, PageData{Active: "blog"})
}

func adminJournalHandler(c *jmux.Context) {
	execTmpl("admin/journal", c, PageData{Active: "journal"})
}

func adminProductsHandler(c *jmux.Context) {
	execTmpl("admin/products", c, PageData{Active: "products"})
}

func adminProductsListHandler(c *jmux.Context) {
	if prodIdStr, ok := c.Params["product_id"]; ok {
		prodId, err := strconv.ParseInt(prodIdStr, 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid product_id")
			return
		}
		prod, err := products.GetProductById(uint64(prodId))
		if err != nil {
			if errors.Is(err, products.ErrNotFound) {
				c.WriteError(http.StatusNotFound, "product not found")
			} else {
				c.WriteError(http.StatusInternalServerError, err.Error())
				log.Printf("error getting product %d: %v", prodId, err)
			}
		} else {
			c.WriteJSON(JsonResp{Data: []products.Product{prod}})
		}
		return
	}
	resp := JsonResp{}
	prods, err := products.GetProducts()
	if err != nil {
		log.Print("error getting products: ", err)
		if prods == nil {
			c.WriteError(http.StatusInternalServerError, err.Error())
			return
		}
		resp.Error = err.Error()
	}
	resp.Data = prods
	c.WriteJSON(resp)
}

func adminProductsEditHandler(c *jmux.Context) {
	ctype := c.ReqHeader().Get("Content-Type")
	prod := products.Product{}
	switch ctype {
	case "application/json":
		if err := c.ReadBodyJSON(&prod); err != nil {
			c.WriteError(http.StatusBadRequest, "bad json")
			return
		}
	case "application/x-www-form-urlencoded":
		if err := c.Request.ParseForm(); err != nil {
			c.WriteError(http.StatusBadRequest, "bad form data")
			return
		}
		id, err := strconv.ParseInt(c.Request.FormValue("id"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid id")
			return
		}
		prod.Id = uint64(id)
		prod.Name = c.Request.FormValue("name")
		prod.Description = c.Request.FormValue("description")
		prod.Webpage = c.Request.FormValue("webpage")
		prod.AppStoreLink = c.Request.FormValue("app-store-link")
		prod.PlayStoreLink = c.Request.FormValue("play-store-link")
		prod.Images = c.Request.Form["images"]
	default:
		c.WriteError(http.StatusBadRequest, "invalid Content-Type")
		return
	}
	if idStr := c.Params["product_id"]; idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid id")
			return
		}
		prod.Id = uint64(id)
	}

	if err := products.EditProduct(prod); err != nil {
		// TODO: check error to see if ID doesn't exist?
		log.Printf("error editing product (%+v): %v", prod, err)
		c.WriteError(http.StatusInternalServerError, err.Error())
		return
	}
}

func adminProductsNewHandler(c *jmux.Context) {
	ctype := c.ReqHeader().Get("Content-Type")
	prod := products.Product{}
	switch ctype {
	case "application/json":
		if err := c.ReadBodyJSON(&prod); err != nil {
			c.WriteError(http.StatusBadRequest, "bad json")
			return
		}
	case "application/x-www-form-urlencoded":
		if err := c.Request.ParseForm(); err != nil {
			c.WriteError(http.StatusBadRequest, "bad form data")
			return
		}
		id, err := strconv.ParseInt(c.Request.FormValue("id"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid id")
			return
		}
		prod.Id = uint64(id)
		prod.Name = c.Request.FormValue("name")
		prod.Description = c.Request.FormValue("description")
		prod.Webpage = c.Request.FormValue("webpage")
		prod.AppStoreLink = c.Request.FormValue("app-store-link")
		prod.PlayStoreLink = c.Request.FormValue("play-store-link")
		prod.Images = c.Request.Form["images"]
		hiddenStr := c.Request.FormValue("hidden")
		prod.Hidden, err = strconv.ParseBool(hiddenStr)
		if hiddenStr != "" && err != nil {
			c.WriteError(http.StatusBadRequest, "invalid hidden")
			return
		}
	default:
		c.WriteError(http.StatusBadRequest, "invalid Content-Type")
		return
	}
	prod, err := products.AddProduct(prod)
	if err != nil {
		// TODO: check error to see if ID doesn't exist?
		log.Printf("error editing product (%+v): %v", prod, err)
		c.WriteError(http.StatusInternalServerError, err.Error())
		return
	}
	c.WriteJSON(JsonResp{Data: prod})
}

func adminProductsIssuesHandler(c *jmux.Context) {
	giq := products.GetIssuesQuery{}
	query := c.Query()

	if prodIdStr := query.Get("product_id"); prodIdStr != "" {
		prodId, err := strconv.ParseInt(prodIdStr, 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "bad product_id")
			return
		}
		giq.FilterProduct = utils.NewT(uint64(prodId))
	}
	if sortBy := query.Get("sort_by"); sortBy != "" {
		giq.SortBy = utils.NewT(sortBy)
	}
	if descStr := query.Get("sort_desc"); descStr != "" {
		desc, err := strconv.ParseBool(descStr)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "bad sort_desc")
			return
		}
		giq.SortDesc = utils.NewT(desc)
	}
	if reason := query.Get("reason"); reason != "" {
		giq.SortBy = utils.NewT(reason)
	}
	if startedStr := query.Get("started"); startedStr != "" {
		started, err := strconv.ParseBool(startedStr)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "bad started")
			return
		}
		giq.FilterStarted = utils.NewT(started)
	}
	if resolvedStr := query.Get("resolved"); resolvedStr != "" {
		resolved, err := strconv.ParseBool(resolvedStr)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "bad resolved")
			return
		}
		giq.FilterResolved = utils.NewT(resolved)
	}

	resp := JsonResp{}
	issues, err := products.GetProductIssues(giq)
	if err != nil {
		if issues == nil {
			log.Printf("error getting issues with params %+v: %v", giq, err)
			c.WriteError(http.StatusInternalServerError, err.Error())
			return
		}
		resp.Error = err.Error()
	}
	resp.Data = issues
	c.WriteJSON(resp)
}

func adminProductsIssueEditHandler(c *jmux.Context) {
	ctype := c.ReqHeader().Get("Content-Type")
	issue := products.ProductIssue{}
	switch ctype {
	case "application/json":
		if err := c.ReadBodyJSON(&issue); err != nil {
			println(err.Error())
			c.WriteError(http.StatusBadRequest, "bad json")
			return
		}
	case "application/x-www-form-urlencoded":
		if err := c.Request.ParseForm(); err != nil {
			c.WriteError(http.StatusBadRequest, "bad form data")
			return
		}
		id, err := strconv.ParseInt(c.Request.FormValue("id"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid id")
			return
		}
		prodId, err := strconv.ParseInt(c.Request.FormValue("product-id"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid product-id")
			return
		}
		issue.ProductId = uint64(prodId)
		issue.Id = uint64(id)
		issue.Email = c.Request.FormValue("email")
		issue.Reason = c.Request.FormValue("reason")
		issue.Subject = c.Request.FormValue("subject")
		issue.Description = c.Request.FormValue("description")
		issue.CreatedAt, err = strconv.ParseInt(c.Request.FormValue("created-at"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid created-at")
			return
		}
		issue.StartedAt, err = strconv.ParseInt(c.Request.FormValue("started-at"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid started-at")
			return
		}
		issue.ResolvedAt, err = strconv.ParseInt(c.Request.FormValue("resolved-at"), 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid resolved-at")
			return
		}
		issue.Ip = c.Request.FormValue("ip")
	default:
		c.WriteError(http.StatusBadRequest, "invalid Content-Type")
		return
	}
	if idStr := c.Params["issue_id"]; idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.WriteError(http.StatusBadRequest, "invalid id")
			return
		}
		issue.Id = uint64(id)
	}

	if err := products.EditProductIssue(issue); err != nil {
		// TODO: check error to see if ID doesn't exist?
		log.Printf("error editing issue (%+v): %v", issue, err)
		c.WriteError(http.StatusInternalServerError, err.Error())
		return
	}
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
		path := cleanPath(c.Path())
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
		if path == "" {
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

const (
	maxLoginAttempts   = 5
	loginAttemptsReset = time.Hour * 12
)

type LoginAttempts struct {
	numAttempts  int
	firstAttempt time.Time
}

var (
	loginAttempts = utils.NewMutex(utils.NewMap[string, LoginAttempts]())
)

func adminLoginHandler(c *jmux.Context) {
	w, r := c.Writer, c.Request

	laMap := *loginAttempts.Lock()
	defer loginAttempts.Unlock()

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	attempts, ok := laMap.GetOk(ip)
	now := time.Now()
	if !ok {
		attempts.firstAttempt = now
	} else if attempts.firstAttempt.Sub(now) >= loginAttemptsReset {
		attempts.numAttempts = 0
	}
	if attempts.numAttempts >= maxLoginAttempts {
		w.WriteHeader(http.StatusUnauthorized)
		execTmpl("admin/login", c, PageData{})
		return
	}

	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if username != adminConfig.Username || password != adminConfig.Password {
		attempts.numAttempts++
		if attempts.numAttempts >= maxLoginAttempts {
			log.Printf("IP %s reached max login attempts (%d)", ip, maxLoginAttempts)
		}
		laMap.Set(ip, attempts)
		w.WriteHeader(http.StatusUnauthorized)
		execTmpl("admin/login", c, PageData{})
		return
	}
	laMap.Delete(ip)
	tokStr, err := newAdminJwt(username)
	if err != nil {
		log.Printf("error creating admin jwt: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("IP %s logged into admin", ip)
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
	// TODO: delete false
	if err != nil || (false && host != remoteIP) {
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

func getTmplPath(tmplName string) string {
	return filepath.Join(tmplsDir, tmplName+".tmpl")
}

func loadTmpl(tmplName string) error {
	var tmpl *template.Template
	var err error
	if strings.HasPrefix(tmplName, "admin/") {
		if tmplName != "admin/login" {
			tmpl, err = template.ParseFiles(
				getTmplPath("admin/base"),
				getTmplPath(tmplName),
			)
		} else {
			tmpl, err = template.ParseFiles(getTmplPath(tmplName))
		}
	} else {
		tmpl, err = template.ParseFiles(baseTmplPath, getTmplPath(tmplName))
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

func cleanPath(path string) string {
	if path != "" && path[0] == '/' {
		return path[1:]
	}
	return path
}

type PageData struct {
	Active string
	Data   any
}

type JsonResp struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func storeSMEntryCreator(name string, priority float32) {
	smUrlEntryCreators.Store(name, makeSMEntryCreator(name, priority))
}

func makeSMEntryCreator(name string, priority float32) SMEntryCreator {
	return func(sitemap.UrlEntry) sitemap.UrlEntry {
		ent := sitemap.UrlEntry{
			Loc:        "https://johnietre.com/" + name,
			Priority:   priority,
			ChangeFreq: sitemap.ChangeFreqMonthly,
		}
		if info, err := os.Stat(getTmplPath(name)); err == nil {
			ent.LastMod = info.ModTime().Format("2006-01-02")
		}
		return ent
	}
}

type SMEntryCreator = func(sitemap.UrlEntry) sitemap.UrlEntry
