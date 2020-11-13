package main

import (
  "log"
  "net"
  "net/http"
  "os"
)

const (
  IP string = "129.119.172.61"
  WEB_PORT string = ":8000"
  CHAT_PORT string = ":8008"
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
  defer LogMessage("Stopping Server")
  LogMessage("Starting Servers")

  webServer := http.Server{
    Addr: IP + WEB_PORT,
    Handler: webRoutes(),
  }
  stocksServer := http.Server{
    Addr: IP + STOCKS_PORT,
    Handler: stocksRoutes(),
  }
  chatServer := http.Server{
    Addr: IP + CHAT_PORT,
    Handler: chatRoutes(),
  }

  go func() {
    log.Panic(webServer.ListenAndServe())
  }()
  go func() {
    log.Panic(chatServer.ListenAndServe())
  }()
  log.Panic(stocksServer.ListenAndServe())
}

// Used to send the most important logs to the logger server
func LogMessage(msg string) {
  // log.Println(msg)
  _, err := logConn.Write([]byte(msg))
  if err != nil {
    log.Println("Error sending log:", err)
  }
}
