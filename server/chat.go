package main

/*
 * Possibly have a chat id as the query
 * Have a special value in the queury for the bot
*/

import (
  "log"
  "net"
  "net/http"
  "os"

  "golang.org/x/net/websocket"
)

const (
  BOT_IP = "localhost"
  BOT_PORT = ":7001"
)

var (
  chatLogger *log.Logger
)

func init() {
  chatLogger = log.New(os.Stdout, "Chat Server: ", log.LstdFlags)
}

func chatRoutes() *http.ServeMux {
  r := http.NewServeMux()
  r.Handle("/", websocket.Handler(u2uHandler))
  r.Handle("/bot", websocket.Handler(u2bHandler))
  return r
}

// u2uHandler handles User-to-User chat
func u2uHandler(ws *websocket.Conn) {
  ws.Close()
}

// u2bHandler handles User-to-Bot chat
func u2bHandler(ws *websocket.Conn) {
  defer ws.Close()
  check := func(err error) bool {
    if err != nil {
      chatLogger.Println(err)
      ws.Write([]byte("ERROR"))
    }
    return err != nil
  }
  conn, err := net.Dial("tcp", BOT_IP+BOT_PORT, )
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