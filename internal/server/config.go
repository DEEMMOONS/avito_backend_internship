package server

import (
  "encoding/json"
  "os"
)

type config struct {
  DBname string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
	Address string `json:"address"`
	Port    string `json:"port"`
}

func CreateConfig(path string) (*config, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }
  config := config{}
  err = json.Unmarshal(data, &config)
  if err != nil {
    return nil, err
  }
  return &config, nil
}
