package main

import (
  _ "fmt"
  _ "net"
)

func main() {
  ln, err := net.Listen("tcp", ":8000")
  if err != nil {
    panic(err)
  }
  for {
    conn, err := ln.Accept()
  }
}
