package main

/*
 * Possibly have a chat id as the query
 * Have a special value in the queury for the bot
*/

import (
  "log"
  "net"
  "os"

  "golang.org/x/net/websocket"
)

const (
  BOT_IP string = "localhost"
  BOT_PORT string = ":7001"
  runSecondHub bool = false
)

var (
  chatLogger *log.Logger
  hub1Chan chan []byte
  hub2Chan chan []byte
)

func init() {
  chatLogger = log.New(os.Stdout, "Chat Server: ", log.LstdFlags)
  hub1Chan = make(chan []byte, 100)
  go startChatHub1()
  if runSecondHub {
    hub2Chan = make(chan []byte, 100)
    go startChatHub2()
  }
}

func chatSocketHandler(ws *websocket.Conn) {
  defer ws.Close()
  for {
    var bmsg [512]byte
    if l, err := ws.Read(bmsg[:]); err != nil {
      if err.Error() == "EOF" {
        return
      }
    } else {
      if runSecondHub {
        select {
        case hub1Chan <- bmsg[:l]:
        case hub2Chan <- bmsg[:l]:
        }
      } else {
        hub1Chan <- bmsg[:l]
      }
    }
  }
  check := func(err error) bool {
    if err != nil {
      chatLogger.Println(err)
      ws.Write([]byte("ERROR"))
    }
    return err != nil
  }
  conn, err := net.Dial("tcp", BOT_IP+BOT_PORT)
  if check(err) {
    return
  }
  _, err = ws.Write([]byte("CONNECTED"))
  if check(err) {
    return
  }
  defer conn.Close()
  var bmsg [512]byte
  for {
    l, err := ws.Read(bmsg[:])
    if err != nil {
      if err.Error() == "EOF" {
        return
      }
      chatLogger.Println(err)
      ws.Write([]byte("ERROR"))
      continue
    }
    _, err = conn.Write(bmsg[:l])
    if check(err) {
      continue
    }
    l, err = conn.Read(bmsg[:])
    if err != nil {
      if err.Error() == "EOF" {
        return
      }
      chatLogger.Println(err)
      ws.Write([]byte("ERROR"))
      continue
    }
    _, err = ws.Write(bmsg[:l])
    check(err)
  }
}

func startChatHub1() {
  for {
    msg := <- hub1Chan
    println(msg)
  }
}

func startChatHub2() {
  for {
    msg := <- hub2Chan
    println(msg)
  }
}