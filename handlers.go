package main

import (
	"net"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "home", Data: githubRepos.Load()}
	if tmpl, loaded := tmpls.Load("home"); !loaded {
		logger.Printf("missing home template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else if err := tmpl.Execute(w, data); err != nil {
		logger.Printf("error executing template: %v", err)
		// NOTE: This would result in probable double-write
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "me"}
	if tmpl, loaded := tmpls.Load("me"); !loaded {
		logger.Printf("missing me template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else if err := tmpl.Execute(w, data); err != nil {
		logger.Printf("error executing template: %v", err)
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "blog"}
	if tmpl, loaded := tmpls.Load("blog"); !loaded {
		logger.Printf("missing blog template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else if err := tmpl.Execute(w, data); err != nil {
		logger.Printf("error executing template: %v", err)
	}
}

func journalHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Active: "journal"}
	if tmpl, loaded := tmpls.Load("journal"); !loaded {
		logger.Printf("missing journal template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	} else if err := tmpl.Execute(w, data); err != nil {
		logger.Printf("error executing template: %v", err)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/parse" {
		if parseTemplate(w, r) {
			w.Write([]byte("Success"))
			return
		}
	}
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
			logger.Printf("error parsing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return false
		}
	}
	logger.Printf("parsed template")
	return true
}

type PageData struct {
	Active string
	Data   any
}
