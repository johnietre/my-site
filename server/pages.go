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
  homeTemplates *template.Template
  stocksTemplates *template.Template
  chatTemplates *template.Template
  convoTemplates *template.Template
  baseTempName string
)

func init() {
  pageLogger = log.New(os.Stdout, "Page Server: ", log.LstdFlags)
  // Conditional possibly unneeded
  if strings.Contains(cwd, "server") {
    baseTempName = "../templates/base.html"
    homeTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/home.html"))
    stocksTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/stocks.html"))
    chatTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/chat.html"))
    convoTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/chat.html"))
  } else {
    baseTempName = "templates/base.html"
    homeTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/home.html"))
    stocksTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/stocks.html"))
    chatTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/chat.html"))
    convoTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/chat.html"))
  }
}

func parse() {
  if strings.Contains(cwd, "server") {
    baseTempName = "../templates/base.html"
    homeTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/home.html"))
    stocksTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/stocks.html"))
    chatTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/chat.html"))
    convoTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "../templates/convo.html"))
  } else {
    baseTempName = "templates/base.html"
    homeTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/home.html"))
    stocksTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/stocks.html"))
    chatTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/chat.html"))
    convoTemplates = template.Must(template.ParseFiles(
      baseTempName,
      "templates/convo.html"))
  }
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
