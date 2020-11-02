package main

import (
  "log"
  "net/http"
  "os"

  "golang.org/x/net/websocket"
)

type Stock struct {
  Sym string `json:"sym"`
}

type Order struct {
  AccountID int64 `json:"account_id"`
  Password string `json:"password"`
  Sym string `json:"sym"`
  Qty string `json:"qty"`
  Side string `json:"side"`
  LimitPrice string `json:"limit_price"`
  StopPrice string `json:"stop_price"`
}

var (
  stocksLogger *log.Logger
)

func init() {
  stocksLogger = log.New(os.Stdout, "Stocks server: ", log.LstdFlags)
}

func stocksRoutes() *http.ServeMux {
  r := http.NewServeMux()
  r.Handle("/", websocket.Handler(symbolHandler))
  r.Handle("/order", websocket.Handler(orderHandler))
  return r
}

func symbolHandler(ws *websocket.Conn) {
  ws.Close()
}

func orderHandler(ws *websocket.Conn) {
  ws.Close()
}