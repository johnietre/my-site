package main

import (
  "log"
  "net"
  "os"
)

var (
  serverLogger *log.Logger
)

func init() {
  serverLogger = log.New(os.Stdout, "DB Server: ", log.LstdFlags)
}

func startServer() {
  ln, err := net.Listen("tcp", IP+":"+PORT)
  if err != nil {
    serverLogger.Panic(err)
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      serverLogger.Println(err)
      continue
    }
    go handle(conn)
  }
}

func handle(conn net.Conn) {
  defer conn.Close()
  return
}