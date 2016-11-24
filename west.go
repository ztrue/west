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
    mt, mreq, err := c.ReadMessage()
    if err != nil {
      log.Println("*** read error:", err)
      break
    }

    go func(mt int, mreq []byte) {
      log.Printf("--- recveived: %s\n", mreq)

      mres, err := proxy.Process(mreq)
      if err != nil {
        log.Println("*** proxy error:", err)
        err = c.WriteJSON(proxy.CometError{Error: err.Error()})
        if err != nil {
          log.Println("*** write err error:", err)
        }
        return
      }

      err = c.WriteMessage(mt, mres)
      if err != nil {
        log.Println("*** write error:", err)
        return
      }

      log.Printf("--- sent: %s\n", mres)
    }(mt, mreq)
  }
}
