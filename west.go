package main

import (
  "encoding/json"
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
  log.SetFlags(0)

  addr := fmt.Sprintf("%s:%s", *host, *port)

  http.HandleFunc("/", handle)
  log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
  c, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println("*** upgrade error:", err)
    return
  }

  log.Println("--- connected")

  defer c.Close()

  for {
    mt, message, err := c.ReadMessage()
    if err != nil {
      log.Println("*** read error:", err)
      break
    }

    log.Printf("--- recveived: %s\n", message)

    var creq = &proxy.CometRequest{}
    err = json.Unmarshal(message, creq)
    if err != nil {
      log.Println("*** parse error:", err)
      break
    }

    cres, err := proxy.Request(creq)
    if err != nil {
      log.Println("*** request error:", err)
      break
    }

    text, err := json.Marshal(cres)
    if err != nil {
      log.Println("*** json error:", err)
      break
    }

    err = c.WriteMessage(mt, text)
    if err != nil {
      log.Println("*** write error:", err)
      break
    }
  }
}
