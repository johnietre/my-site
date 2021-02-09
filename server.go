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
		ip = "192.168.1.137"
	}
	if port == "" {
		port = "443"
	}
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", static))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

  log.Printf("Running on %s:%s", ip, port)
	log.Panic(http.ListenAndServe(ip+":"+port, nil))
}
