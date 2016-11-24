package proxy

import (
  "encoding/json"
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

func DecodeRequest(message []byte) (*CometRequest, error) {
  var creq = &CometRequest{}
  err := json.Unmarshal(message, creq)
  if err != nil {
    return nil, err
  }
  return creq, nil
}

func ConvertResponse(res *http.Response, id string) (*CometResponse, error) {
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
  return &CometResponse{id, res.StatusCode, res.Header, string(body)}, nil
}

func EncodeResponse(cres *CometResponse) ([]byte, error) {
  return json.Marshal(cres)
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

func Process(message []byte) ([]byte, error) {
  creq, err := DecodeRequest(message)
  if err != nil {
    return nil, err
  }

  cres, err := Request(creq)
  if err != nil {
    return nil, err
  }

  return EncodeResponse(cres)
}
