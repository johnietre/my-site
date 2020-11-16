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
  IP string = "localhost"
  PORT string = ":8000"
  STOCKS_PORT string = ":8080"
)

var (
  logConn net.Conn
  cwd string
)

func init() {
  var err error
  cwd, err = os.Getwd()
  if err != nil {
    log.Panic(err)
  }
  logConn, err = net.Dial("tcp", "localhost:7000")
  if err != nil {
    log.Println("Error connecting to logger:", err)
  }
}

func main() {
  // defer LogMessage("Stopping Server")
  // LogMessage("Starting Servers")

  server := http.Server{
    Addr: IP + PORT,
    Handler: routes(),
  }
  log.Panic(server.ListenAndServe())
}

func routes() *http.ServeMux {
  r := http.NewServeMux()
  r.HandleFunc("/", homePageHandler)
  r.HandleFunc("/chat/", chatPageHandler)
  r.HandleFunc("/stocks", stocksPageHandler)
  r.Handle("/chatsocket/", websocket.Handler(chatSocketHandler))

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
  log.Println(msg)
  // _, err := logConn.Write([]byte(msg))
  // if err != nil {
  //   log.Println("Error sending log:", err)
  // }
}
