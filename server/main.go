package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	DEFAULT_IP          string = "192.168.1.137"
	DEFAULT_PORT        string = "443"
	DEFAULT_STOCKS_PORT string = "8080"
	DEFAULT_LOG_IP      string = "localhost"
	DEFAULT_LOG_PORT    string = "7000"
)

var (
	IP             string = os.Getenv("IP")          // IP for the server
	PORT           string = os.Getenv("PORT")        // port for the server
	STOCKS_PORT    string = os.Getenv("STOCKS_PORT") // port for the stocks socket
	LOG_IP         string = os.Getenv("LOG_IP")      // IP for the logging server
	LOG_PORT       string = os.Getenv("LOG_PORT")    // port for the logging server
	logConn        net.Conn
	connectedToLog bool = true
	cwd            string
)

func init() {
	var err error
	// Get the current directory; based on that info, point to the correct
	// place for the templates and static files
	cwd, err = os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	// Check to make sure all the IPs and PORTs have been set
	if IP == "" {
		log.Printf(`Environ variable "IP" not set... using "%s"`, DEFAULT_IP)
		log.Println("")
		IP = DEFAULT_IP
	}
	if PORT == "" {
		log.Printf(`Environ variable "PORT" not set... using "%s"`, DEFAULT_PORT)
		log.Println("")
		PORT = DEFAULT_PORT
	}
	if STOCKS_PORT == "" {
		log.Printf(`Environ variable "STOCKS_PORT" not set... using "%s"`, DEFAULT_STOCKS_PORT)
		log.Println("")
		STOCKS_PORT = DEFAULT_STOCKS_PORT
	}
	if LOG_IP == "" {
		log.Printf(`Environ variable "LOG_IP" not set... using "%s"`, DEFAULT_LOG_IP)
		log.Println("")
		LOG_IP = DEFAULT_LOG_IP
	}
	if LOG_PORT == "" {
		log.Printf(`Environ variable "LOG_PORT" not set... using "%s"`, DEFAULT_LOG_PORT)
		log.Println("")
		LOG_PORT = DEFAULT_LOG_PORT
	}

	// Connect to the log server
	logConn, err = net.Dial("tcp", LOG_IP+":"+LOG_PORT)
	if err != nil {
		log.Println("Error connecting main to logger:", err)
		connectedToLog = false
	}
}

func main() {
	defer LogMessage("Stopping Server")
	LogMessage("Starting Servers")

	server := http.Server{
		Addr:    IP + ":" + PORT,
		Handler: routes(),
	}
	log.Panic(server.ListenAndServe())
}

func routes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/", homePageHandler)
	r.HandleFunc("/chat/", chatPageHandler)
	r.HandleFunc("/stocks", stocksPageHandler)
	r.HandleFunc("/login", loginPageHandler)
	r.HandleFunc("/register", registerPageHandler)
	r.Handle("/chatsocket/", websocket.Handler(chatSocketHandler))

	// Check the cwd to see where the templates and static directories are
	var staticPath string
	if strings.Contains(cwd, "server") {
		staticPath = "../static"
	} else {
		staticPath = "static"
	}
	static := http.FileServer(http.Dir(staticPath))
	r.Handle("/static/", http.StripPrefix("/static", static))

	return r
}

// LogMessage is used to send the most important logs to the logger server
func LogMessage(msg string) {
	if !connectedToLog {
		log.Println(msg)
		return
	}
	_, err := logConn.Write([]byte(msg))
	if err != nil {
		log.Println("Error sending log:", err)
	}
}
