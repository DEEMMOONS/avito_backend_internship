package main

import (
    "log"
    "github.com/DEEMMOONS/avito_backend_internship/tree/develop/internal/server"
)

func main() {
  server, err := server.NewServer("app/configs/config.json")
  if err != nil {
    log.Fatal(err)
  }
  server.Start()
}
