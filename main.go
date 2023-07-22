package main

import (
	"context"
  "errors"
	"flag"
	"html/template"
  "io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
  "sync/atomic"
	"syscall"
	"time"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.LUTC)
  remoteIP string
  indexPath = "./index.html"
  tmpVal atomic.Value
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8000", "Address to run on")
  flag.StringVar(&remoteIP, "remote-ip", "", "Remote IP to check for to parse")
	certPath := flag.String("cert", "", "Path to cert file")
	keyPath := flag.String("key", "", "Path to key file")
	logPath := flag.String("log", "", "Path to log file (empty routes to stderr")
	flag.StringVar(&indexPath, "index", "./index.html", "Path to index file")
	staticPath := flag.String("static", "./static", "Path to static files directory")
	flag.Parse()

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

	router := http.NewServeMux()

	static := http.FileServer(http.Dir(*staticPath))
	router.Handle("/static/", http.StripPrefix("/static", static))

	indexTmp, err := template.ParseFiles(indexPath)
	if err != nil {
		logger.Fatalf("error parsing index file: %v", err)
	}
  // Check to make sure the template executes without err
  if err := indexTmp.Execute(io.Discard, nil); err != nil {
    logger.Fatalf("error executing template: %v", err)
  }
  tmpVal.Store(indexTmp)

  router.HandleFunc("/", handler)

	srvr := &http.Server{
		Addr:              *addr,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
		ErrorLog:          logger,
	}

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

func handler(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  if path == "/parse" {
    if parseTemplate(w, r) {
      return
    }
  }
  if err := tmpVal.Load().(*template.Template).Execute(w, nil); err != nil {
    //logger.Printf("error executing template: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}

// Returns true if successful
func parseTemplate(w http.ResponseWriter, r *http.Request) bool {
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil || host != remoteIP {
      return false
    }
    tmp, err := template.ParseFiles(indexPath)
    if err != nil {
      logger.Printf("error parsing template: %v", err)
      http.Error(w, "Internal Server Error", http.StatusInternalServerError)
      return false
    }
    if err := tmp.Execute(w, nil); err != nil {
      logger.Printf("error executing template: %v", err)
      http.Error(w, "Internal Server Error", http.StatusInternalServerError)
      return false
    }
    tmpVal.Store(tmp)
    logger.Printf("parsed template")
    return true
}
