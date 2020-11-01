package main

import (
  "html/template"
  "log"
  "net/http"
  "os"
)

var (
  pageLogger *log.Logger
)

func StartPageServer() {
  pageLogger = log.New(os.Stdout, "Page Server: ", log.LstdFlags)
  server := http.Server{
    Addr: IP + PAGE_PORT,
    Handler: pageRoutes(),
  }
  err := server.ListenAndServe()
  LogMessage("Stopping server from pageServer.go")
  pageLogger.Panic(err)
}

func pageRoutes() *http.ServeMux {
  r := http.NewServeMux()
  r.HandleFunc("/", homePageHandler)

  static := http.FileServer(http.Dir("../static"))
  r.Handle("/static/", http.StripPrefix("/static", static))

  return r
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
  if ts, err := template.ParseFiles("../templates/home.html"); err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    pageLogger.Println(err)
    return
  } else {
    ts.Execute(w, nil)
  }
}