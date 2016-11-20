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

type CometRequest struct {
  Id string `json:"id"`
  Method string `json:"method"`
  Url string `json:"url"`
}

type CometResponse struct {
  Id string `json:"id"`
  StatusCode int `json:"status"`
  Header http.Header `json:"headers"`
  Body string `json:"body"`
}

func Convert(res *http.Response, id string) (*CometResponse, error) {
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  return &CometResponse{id, res.StatusCode, res.Header, string(body)}, nil
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

    var creq = &CometRequest{}
    err = json.Unmarshal(message, creq)
    if err != nil {
      log.Println("*** parse error:", err)
      break
    }

    cres, err := request(creq)
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

func request(creq *CometRequest) (*CometResponse, error) {
  req, err := http.NewRequest(creq.Method, creq.Url, nil)
  if err != nil {
    return nil, err
  }

  client := &http.Client{}
  res, err := client.Do(req)
  if res != nil {
    defer res.Body.Close()
  }
  if err != nil {
    return nil, err
  }

  log.Println(res)

  return Convert(res, creq.Id)
}
