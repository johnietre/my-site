package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/websocket"
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
		log.Println(`Environ variable "IP" not set... using "129.119.172.61"`)
		IP = "129.119.172.61"
	}
	if PORT == "" {
		log.Println(`Environ variable "PORT" not set... using "8000"`)
		PORT = "8000"
	}
	if STOCKS_PORT == "" {
		log.Println(`Environ variable "STOCKS_PORT" not set... using "8080"`)
		STOCKS_PORT = "8080"
	}
	if LOG_IP == "" {
		log.Println(`Environ variable "LOG_IP" not set... using "localhost"`)
		LOG_IP = "localhost"
	}
	if LOG_PORT == "" {
		log.Println(`Environ variable "LOG_PORT" not set... using "7000"`)
		LOG_PORT = "7000"
	}

	// Connect to the log server
	logConn, err = net.Dial("tcp", LOG_PORT+":"+LOG_PORT)
	if err != nil {
		log.Println("Error connecting to logger:", err)
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

// Used to send the most important logs to the logger server
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
