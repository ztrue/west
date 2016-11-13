package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "log"
  "io/ioutil"
  "net/http"
  "github.com/gorilla/websocket"
)

var host = flag.String("host", "localhost", "http service host")
var port = flag.String("port", "5000", "http service port")

var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

type CometResponse struct  {
  StatusCode int `json:"status"`
  Header http.Header `json:"headers"`
  Body string `json:"body"`
}

func Convert(res *http.Response) (*CometResponse, error) {
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  return &CometResponse{res.StatusCode, res.Header, string(body)}, nil
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

    res, err := request(string(message))
    if err != nil {
      log.Println("*** request error:", err)
      break
    }

    log.Println(res)

    cr, err := Convert(res)
    if err != nil {
      log.Println("*** convert error:", err)
      break
    }

    text, err := json.Marshal(cr)
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

func request(url string) (*http.Response, error) {
  method := "GET"
  req, err := http.NewRequest(method, url, nil)
  if err != nil {
    return nil, err
  }

  client := &http.Client{}
  return client.Do(req)
}
