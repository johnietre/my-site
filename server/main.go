package main

import (
  "log"
  "net"
)

const (
  IP string = "129.119.172.180"
  PAGE_PORT string = ":8000"
)

var (
  conn net.Conn
)

func main() {
  var err error
  conn, err = net.Dial("tcp", "localhost:50989")
  if err != nil {
    log.Println("Error connecting to logger:", err)
  }
  defer LogMessage("Stopping Server")
  LogMessage("Starting Server")
  StartPageServer()
}

// Used to send the most important logs to the logger server
func LogMessage(msg string) {
  _, err := conn.Write([]byte(msg))
  if err != nil {
    log.Println("Error sending log:", err)
  }
}