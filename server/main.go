package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/johnietre/my-site/server/apps"
	"github.com/johnietre/my-site/server/blogs"
	"github.com/johnietre/my-site/server/handlers"
	"github.com/johnietre/my-site/server/repos"
)

var (
//baseDir = "./"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func main() {
	baseDir := flag.String(
		"base-dir", ".",
		"Path to base directory where all other directories are located",
	)

	addr := flag.String("addr", "127.0.0.1:8000", "Address to run on")
	staticDir := flag.String("static-dir", "$(base-dir)/static", "Path to static directory")
	adminConfigPath := flag.String(
		"admin-config",
		"$(base-dir)/server/admin-config.json",
		"Path to admin config file",
	)
	newAdminConfig := flag.Bool(
		"new-admin-config",
		false,
		"Create admin config file",
	)
	certPath := flag.String("cert", "", "Path to cert file")
	keyPath := flag.String("key", "", "Path to key file")

	logPath := flag.String("log", "", "Path to log file (empty routes to stderr")
	appsDbPath := flag.String(
		"apps-db", "$(base-dir)/apps/apps.db", "Path to apps database",
	)
	blogsDir := flag.String(
		"blogs-dir", "$(base-dir)/blogs", "Path to blogs directory",
	)
	blogsDbPath := flag.String(
		"blogs-db", "$(blogs-dir)/blogs.db", "Path to blogs database",
	)
	tmplsDir := flag.String(
		"tmpls-dir", "$(base-dir)/templates", "Path to templates directory",
	)
	remoteIP := flag.String("remote-ip", "", "Remote IP to check for to parse")

	flag.Parse()

	if *newAdminConfig {
		f, err := os.Create("admin-config.json")
		if err != nil {
			log.Fatalf("error creating admin-config.json: %v", err)
		}
		defer f.Close()
		e := json.NewEncoder(f)
		e.SetIndent("", "  ")
		if err := e.Encode(handlers.AdminConfig{}); err != nil {
			log.Fatalf("error creating admin-config.json: %v", err)
		}
		return
	}

	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
		log.SetOutput(f)
	}

	if *baseDir == "" {
		*baseDir = "./"
	}
	*baseDir = filepath.Join(*baseDir, "a")
	// Add a dummy part to the path to add a trailing slug
	*baseDir = (*baseDir)[:len(*baseDir)-1]

	replacePath(adminConfigPath, "$(base-dir)/", *baseDir)
	replacePath(staticDir, "$(base-dir)/", *baseDir)
	replacePath(appsDbPath, "$(base-dir)/", *baseDir)
	replacePath(blogsDir, "$(base-dir)/", *baseDir)
	replacePath(blogsDbPath, "$(blogs-dir)/", *blogsDir)
	replacePath(tmplsDir, "$(base-dir)/", *baseDir)

	adminConfig := handlers.AdminConfig{}
	f, err := os.Open(*adminConfigPath)
	if err != nil {
		log.Fatalf("error opening admin config file: %v", err)
	}
	err = json.NewDecoder(f).Decode(&adminConfig)
	f.Close()
	if err != nil {
		log.Fatalf("error decoding admin config file: %v", err)
	}

	if err := handlers.InitHandlers(*tmplsDir, *remoteIP, adminConfig); err != nil {
		log.Fatalf("error initializing handlers: %v", err)
	} else if err = apps.InitApps(*appsDbPath); err != nil {
		log.Fatalf("error initializing apps: %v", err)
	} else if err = blogs.InitBlogs(*blogsDir, *blogsDbPath); err != nil {
		log.Fatalf("error initializing blogs: %v", err)
	} else if err = repos.InitRepos(); err != nil {
		log.Fatalf("error initializing blogs: %v", err)
	}

	doneChan := make(chan bool, 2)
	err = RunServer(*addr, *staticDir, *keyPath, *certPath, doneChan)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error running server: %v", err)
	}
	<-doneChan
	log.Print("finished")
}

func RunServer(addr, staticDir, keyPath, certPath string, doneChan chan<- bool) error {
	lc := net.ListenConfig{
		Control: func(network, addr string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if err := syscall.SetsockoptInt(
					int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1,
				); err != nil {
					log.Fatalf("error setting SO_REUSEADDR: %v", err)
				}
			})
		},
	}
	ls, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalf("error starting listener: %v", err)
	}

	srvr := &http.Server{
		Addr:              addr,
		Handler:           handlers.CreateRouter(staticDir),
		ReadHeaderTimeout: time.Second * 10,
	}

	interruptChan := make(chan os.Signal, 5)
	go func() {
		<-interruptChan
		log.Print("shutting down")
		go func() {
			srvr.Shutdown(context.Background())
			doneChan <- true
		}()
		<-interruptChan
		log.Print("closing")
		srvr.Close()
		doneChan <- true
	}()
	signal.Notify(interruptChan, os.Interrupt)

	log.Printf("Running on %s", srvr.Addr)
	if certPath != "" || keyPath != "" {
		err = srvr.ServeTLS(ls, certPath, keyPath)
	} else {
		err = srvr.Serve(ls)
	}
	return err
}

func replacePath(strPtr *string, from, to string) {
	if strings.HasPrefix(*strPtr, from) {
		*strPtr = filepath.Clean(strings.Replace(*strPtr, from, to, 1))
	}
}
