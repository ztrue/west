package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/websocket"
  "./proxy"
)

var host = flag.String("host", "localhost", "http service host")
var port = flag.String("port", "5000", "http service port")

var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func main() {
  flag.Parse()
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)

  addr := fmt.Sprintf("%s:%s", *host, *port)

  http.HandleFunc("/", handle)
  log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
  c, err := upgrader.Upgrade(w, r, nil)
  if c != nil {
    defer c.Close()
  }
  if err != nil {
    log.Println("*** upgrade error:", err)
    return
  }

  log.Println("--- connected")

  for {
    var creq = &proxy.CometRequest{}
    if err := c.ReadJSON(creq); err != nil {
      log.Println("*** read error:", err)
      break
    }

    go func(creq *proxy.CometRequest) {
      log.Println("--- recveived:", creq)

      var message interface{}
      message, err := proxy.Request(creq)
      if err != nil {
        message = proxy.ConvertError(err, "")
      }

      if err := c.WriteJSON(message); err != nil {
        log.Println("*** write error:", err)
      } else {
        log.Println("--- sent:", message)
      }
    }(creq)
  }
}
