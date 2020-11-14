package main

import (
  "log"
  "net"
  "os"
)

const (
  IP string = "192.168.1.146"
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

  go startWeb()
  go startChat()
  startStocks()
}

// Used to send the most important logs to the logger server
func LogMessage(msg string) {
  // log.Println(msg)
  _, err := logConn.Write([]byte(msg))
  if err != nil {
    log.Println("Error sending log:", err)
  }
}
