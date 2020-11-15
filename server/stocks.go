package main
/*
import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"

  "golang.org/x/net/websocket"
)

// Stock is used to hold stock ticker information
type Stock struct {
  Sym string `json:"sym"`
}

// Order is used to hold all relavent order info
type Order struct {
  AccountID int64 `json:"account_id"`
  Password string `json:"password"`
  Sym string `json:"sym"`
  Qty int `json:"qty"`
  Side string `json:"side"`
  LimitPrice float64 `json:"limit_price"`
  StopPrice float64 `json:"stop_price"`
}

var (
  stocksLogger *log.Logger
)

func init() {
  stocksLogger = log.New(os.Stdout, "Stocks server: ", log.LstdFlags)
}

func startStocks() {
  stocksServer := http.Server{
    Addr: IP + STOCKS_PORT,
    Handler: stocksRoutes(),
  }
  log.Panic(stocksServer.ListenAndServe())
}

func stocksRoutes() *http.ServeMux {
  r := http.NewServeMux()
  r.Handle("/", websocket.Handler(symbolHandler))
  r.Handle("/order", websocket.Handler(orderHandler))
  return r
}

func symbolHandler(ws *websocket.Conn) {
  websocket.Message.Send(ws, "Not implemented")
  ws.Close()
}

func orderHandler(ws *websocket.Conn) {
  defer ws.Close()
  var order Order
  var bmsg [512]byte
  for {
    l, err := ws.Read(bmsg[:])
    if err != nil {
      if err.Error() == "EOF" {
        return
      }
      return
    }
    if err := json.Unmarshal(bmsg[:l], &order); err != nil {
      continue
    }
    fmt.Println("Order", order)
  }
}
*/