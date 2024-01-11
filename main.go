package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	utils "github.com/johnietre/utils/go"
)

var (
	logger   = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.LUTC)
	remoteIP string

	baseTmplPath string
	tmplsPath    = "./templates"
	tmplNames    = []string{"home", "me", "blog", "journal"}
	tmpls        = utils.NewSyncMap[string, *template.Template]()

	githubRepos = utils.NewAValue[[]RepoInfo]([]RepoInfo{})
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8000", "Address to run on")
	flag.StringVar(&remoteIP, "remote-ip", "", "Remote IP to check for to parse")
	certPath := flag.String("cert", "", "Path to cert file")
	keyPath := flag.String("key", "", "Path to key file")
	logPath := flag.String("log", "", "Path to log file (empty routes to stderr")
	flag.StringVar(&tmplsPath, "tmpls", "./templates", "Path to templates directory")
	staticPath := flag.String("static", "./static", "Path to static files directory")
	flag.Parse()

	baseTmplPath = filepath.Join(tmplsPath, "base.tmpl")

	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logger.Fatalf("error opening log file: %v", err)
		}
		logger.SetOutput(f)
	}

	lc := net.ListenConfig{
		Control: func(network, addr string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if err := syscall.SetsockoptInt(
					int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1,
				); err != nil {
					logger.Fatalf("error setting SO_REUSEADDR: %v", err)
				}
			})
		},
	}
	ls, err := lc.Listen(context.Background(), "tcp", *addr)
	if err != nil {
		logger.Fatalf("error starting listener: %v", err)
	}

	// Populate the routes
	router := http.NewServeMux()

	static := http.FileServer(http.Dir(*staticPath))
	router.Handle("/static/", http.StripPrefix("/static", static))

	for _, name := range tmplNames {
		if err := loadTmpl(name); err != nil {
			logger.Fatal(err)
		}
	}
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/me", meHandler)
	router.HandleFunc("/blog", blogHandler)
	router.HandleFunc("/journal", journalHandler)

	router.Handle("/admin/", http.StripPrefix("/admin", http.HandlerFunc(adminHandler)))

	srvr := &http.Server{
		Addr:              *addr,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
		ErrorLog:          logger,
	}

	// Load the GitHub repos
	if err := refreshRepos(); err != nil {
		logger.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(time.Minute)
			if err := refreshRepos(); err != nil {
				logger.Println(err)
			}
		}
	}()

	interruptChan, doneChan := make(chan os.Signal, 5), make(chan bool, 2)
	go func() {
		<-interruptChan
		logger.Print("shutting down")
		go func() {
			srvr.Shutdown(context.Background())
			doneChan <- true
		}()
		<-interruptChan
		logger.Print("closing")
		srvr.Close()
		doneChan <- true
	}()
	signal.Notify(interruptChan, os.Interrupt)

	logger.Printf("Running on %s", srvr.Addr)
	if *certPath != "" || *keyPath != "" {
		err = srvr.ServeTLS(ls, *certPath, *keyPath)
	} else {
		err = srvr.Serve(ls)
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("error running server: %v", err)
	}
	<-doneChan
	logger.Print("finished")
}

func loadTmpl(tmplName string) error {
	tmpl, err := template.ParseFiles(
		baseTmplPath, filepath.Join(tmplsPath, tmplName+".tmpl"),
	)
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

func refreshRepos() error {
	resp, err := http.Get(
		"https://api.github.com/users/johnietre/repos?sort=pushed&direction=desc&per_page=3",
	)
	if err != nil {
		return fmt.Errorf("Error getting repos json: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Received non-200 status: %s", resp.Status)
	}
	reposMaps := []RepoInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&reposMaps); err != nil {
		return fmt.Errorf("Error decoding repos json: %v", err)
	}
	repos := make([]RepoInfo, 0, len(reposMaps))
	for _, repo := range reposMaps {
		repos = append(repos, repo)
		/*
		   iUrl, ok := repo["html_url"]
		   if !ok {
		     logger.Println("Missing repo html_url")
		   }
		   url, ok := iUrl.(string)
		   if !ok {
		     logger.Println("Invalid repo html_url received, got %v", iUrl)
		   }
		   iName, ok := repo["name"]
		   if !ok {
		     logger.Println("Missing repo name")
		   }
		   name, ok := iName.(string)
		   if !ok {
		     logger.Println("Invalid repo name received, got %v", iName)
		   }
		*/
	}
	githubRepos.Store(repos)
	return nil
}

type RepoInfo struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}
