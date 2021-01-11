package main

import (
	"log"
	"net/http"
)

func main() {
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", static))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Panic(http.ListenAndServe("192.168.1.125:8008", nil))
}
