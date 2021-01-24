package main

import (
	"log"
	"net/http"
  "os"
)

func main() {
  ip := os.Getenv("IP")
  port := os.Getenv("PORT")
  if ip == "" {
    ip = "localhost"
  }
  if port == "" {
    port = "8000"
  }
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", static))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

  log.Panic(http.ListenAndServe(ip + ":" + port, nil))
}
