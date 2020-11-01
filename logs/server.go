package main

import (
  "log"
  "net"
  "os"
  "strconv"
)

const (
  IP string = "localhost"
  PORT string = "50989"
)

func main() {
  file, err := os.OpenFile("main.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  if err != nil {
    log.Fatalln("Logger Error:", err)
  }
  log.SetOutput(file)

  ln, err := net.Listen("tcp", IP+PORT)
  if err != nil {
    log.Fatalln("Error setting up logger:", err)
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Println(err)
      continue
    }
    go handle(conn)
  }
}

func handle(conn net.Conn) {
  defer conn.Close()
  var bmsg [128]byte
  conn.Write([]byte("9CONNECTED"))
  var l int
  var err error
  for {
    conn.Read(bmsg[:3])
    l, err = strconv.Atoi(string(bmsg[:3]))
    if err != nil {
      log.Println("Logger Error:", err)
      continue
    }
    conn.Read(bmsg[:l])
    log.Println(string(bmsg[:]))
  }
}