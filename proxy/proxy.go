package proxy

import (
  "io/ioutil"
  "net/http"
)

type CometRequest struct {
  Id string `json:"id"`
  Method string `json:"method"`
  Url string `json:"url"`
  Header http.Header `json:"headers"`
}

type CometResponse struct {
  Id string `json:"id"`
  StatusCode int `json:"status"`
  Header http.Header `json:"headers"`
  Body string `json:"body"`
}

type CometError struct {
  Id string `json:"id"`
  Error string `json:"error"`
}

func ConvertResponse(res *http.Response, id string) (*CometResponse, error) {
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
  return &CometResponse{id, res.StatusCode, res.Header, string(body)}, nil
}

func Request(creq *CometRequest) (*CometResponse, error) {
  req, err := http.NewRequest(creq.Method, creq.Url, nil)
  if err != nil {
    return nil, err
  }

  req.Header = creq.Header

  client := &http.Client{}
  res, err := client.Do(req)
  if res != nil {
    defer res.Body.Close()
  }
  if err != nil {
    return nil, err
  }

  return ConvertResponse(res, creq.Id)
}
