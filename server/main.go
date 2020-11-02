package main

import (
  "log"
  // "net"
  "net/http"
  "os"
)

const (
  IP string = "129.119.172.61"
  WEB_PORT string = ":8000"
  CHAT_PORT string = ":8008"
)

var (
  // logCon net.Conn
  cwd string
)

func init() {
  var err error
  cwd, err = os.Getwd()
  if err != nil {
    log.Panic(err)
  }
}

func main() {
  // var err error
  // logConn, err = net.Dial("tcp", "localhost:7000")
  // if err != nil {
  //   log.Println("Error connecting to logger:", err)
  // }
  // defer LogMessage("Stopping Server")
  // LogMessage("Starting Servers")

  webServer := http.Server{
    Addr: IP + WEB_PORT,
    Handler: webRoutes(),
  }
  chatServer := http.Server{
    Addr: IP + CHAT_PORT,
    Handler: chatRoutes(),
  }

  go webServer.ListenAndServe()
  log.Panic(chatServer.ListenAndServe())
}

// Used to send the most important logs to the logger server
func LogMessage(msg string) {
  log.Println(msg)
  // _, err := logConn.Write([]byte(msg))
  // if err != nil {
  //   log.Println("Error sending log:", err)
  // }
}