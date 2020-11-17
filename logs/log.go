package main

import (
  "log"
  "net"
  "os"
  "strconv"
  "strings"
)

var (
  IP string = os.Getenv("LOG_IP")
  PORT string = os.Getenv("LOG_PORT")
)

func main() {
  // Get the current working directory (either MySite or logs)
  // Based on the cwd, make sure the log file is being placed in the correct place
  var path string
  cwd, err := os.Getwd()
  if err != nil {
    log.Panic(err)
  }
  if strings.Contains(cwd, "logs") {
    path = "./"
  } else {
    path = "./logs/"
  }
  file, err := os.OpenFile(path+"main.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  if err != nil {
    log.Fatalln("Logger Error:", err)
  }
  log.SetOutput(file)

  // Check to make sure the IP and PORT environment variables have been set
  if IP == "" {
    log.Println(`Environ variable "LOG_IP" not set... using "localhost"`)
    IP = "localhost"
  }
  if PORT == "" {
    log.Println(`Environ variable "LOG_PORT" not set... using "7000"`)
  }
  // Start listening to logs
  ln, err := net.Listen("tcp", IP+":"+PORT)
  if err != nil {
    log.Fatalln("Error setting up logger:", err)
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Println(err)
      continue
    }
    go handleLogConn(conn)
  }
}

// Handle each program's logging
func handleLogConn(conn net.Conn) {
  defer conn.Close()
  var bmsg [128]byte
  conn.Write([]byte("9CONNECTED"))
  var l int
  var err error
  for {
    conn.Read(bmsg[:3])
    l, err = strconv.Atoi(string(bmsg[:3]))
    if err != nil {
      // log.Println("Logger Error:", err)
      continue
    }
    conn.Read(bmsg[:l])
    // log.Println(string(bmsg[:]))
  }
}