package main

import (
  "html/template"
  "log"
  "net/http"
  "os"
  "strings"
)

/*
 Handle the GET requests to get the specific convo
*/

var (
  pageLogger *log.Logger
  templatesPath string
  staticPath string
)

func init() {
  pageLogger = log.New(os.Stdout, "Page Server: ", log.LstdFlags)
  if strings.Contains(cwd, "server") {
    templatesPath = "../templates/"
    staticPath = "../static/"
  } else {
    templatesPath = "templates/"
    staticPath = "static/"
  }
}

func startWeb() {
  webServer := http.Server{
    Addr: IP + WEB_PORT,
    Handler: webRoutes(),
  }
  log.Panic(webServer.ListenAndServe())
}

func webRoutes() *http.ServeMux {
  r := http.NewServeMux()
  r.HandleFunc("/", homePageHandler)
  r.HandleFunc("/stocks", stocksPageHandler)
  r.HandleFunc("/chat", chatPageHandler)
  r.HandleFunc("/convo", convoPageHandler)

  static := http.FileServer(http.Dir(staticPath))
  r.Handle("/static/", http.StripPrefix("/static", static))

  return r
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
  if ts, err := template.ParseFiles(templatesPath+"home.html"); err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    pageLogger.Println(err)
    return
  } else {
    ts.Execute(w, nil)
  }
}

func stocksPageHandler(w http.ResponseWriter, r *http.Request) {
  if ts, err := template.ParseFiles(templatesPath+"stocks.html"); err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    pageLogger.Println(err)
    return
  } else {
    ts.Execute(w, nil)
  }
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
  if ts, err := template.ParseFiles(templatesPath+"chat.html"); err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    pageLogger.Println(err)
    return
  } else {
    ts.Execute(w, nil)
  }
}

func convoPageHandler(w http.ResponseWriter, r *http.Request) {
  if ts, err := template.ParseFiles(templatesPath+"convo.html"); err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    pageLogger.Println(err)
    return
  } else {
    ts.Execute(w, nil)
  }
}
