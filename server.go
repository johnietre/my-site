package main

import (
  "flag"
	"log"
	"net/http"
)

func main() {
  addr := flag.String("addr", "127.0.0.1:8000", "Address to run on")
  flag.Parse()

	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", static))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

  log.Printf("Running on %s", *addr)
	log.Panic(http.ListenAndServe(*addr, nil))
}
